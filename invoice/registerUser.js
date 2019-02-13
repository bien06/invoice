'use strict';
/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */
/*
 * Register and Enroll a user
 */

var Fabric_Client = require('fabric-client');
var Fabric_CA_Client = require('fabric-ca-client');

var path = require('path');
var util = require('util');
var os = require('os');

//
var fabric_client = new Fabric_Client();
var fabric_ca_client = null;
var admin_user = null;

var member_user_samsung = null;
var member_user_xiaomi = null;
var member_user_unionbank = null;

var secret_samsung = null;
var secret_xiaomi = null;
var secret_unionbank = null;

var store_path = path.join(__dirname, 'hfc-key-store');
console.log(' Store path:' + store_path);

// create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
Fabric_Client.newDefaultKeyValueStore({
    path: store_path
}).then((state_store) => {
    // assign the store to the fabric client
    fabric_client.setStateStore(state_store);
    var crypto_suite = Fabric_Client.newCryptoSuite();
    // use the same location for the state store (where the users' certificate are kept)
    // and the crypto store (where the users' keys are kept)
    var crypto_store = Fabric_Client.newCryptoKeyStore({
        path: store_path
    });
    crypto_suite.setCryptoKeyStore(crypto_store);
    fabric_client.setCryptoSuite(crypto_suite);
    var tlsOptions = {
        trustedRoots: [],
        verify: false
    };
    // be sure to change the http to https when the CA is running TLS enabled
    fabric_ca_client = new Fabric_CA_Client('http://localhost:7054', null, '', crypto_suite);

    // first check to see if the admin is already enrolled
    return fabric_client.getUserContext('admin', true);
}).then((user_from_store) => {
    if (user_from_store && user_from_store.isEnrolled()) {
        console.log('Successfully loaded admin from persistence');
        admin_user = user_from_store;
    } else {
        throw new Error('Failed to get admin.... run enrollAdmin.js');
    }

    // at this point we should have the admin user
    // first need to register the user with the CA server
    //var attributes = {username:"Amol:ecert",org:"ICICI:ecert"};
    let attributes = [{
            name: "username",
            value: "Samsung",
            ecert: true
        }, // Supplier
        {
            name: "username",
            value: "Xiaomi",
            ecert: true
        }, // OEM
        {
            name: "username",
            value: "UnionBank",
            ecert: true
        } // Bank

    ];

    return fabric_ca_client
        .register({
            enrollmentID: 'Samsung',
            affiliation: 'org1.department1',
            role: 'supplier',
            attrs: attributes
        }, admin_user)
        .then((samsung_secret) => {
            secret_samsung = samsung_secret
            return fabric_ca_client
                .register({
                    enrollmentID: 'Xiaomi',
                    affiliation: 'org1.department1',
                    role: 'oem',
                    attrs: attributes
                }, admin_user)
                .then((xiaomi_secret) => {
                    secret_xiaomi = xiaomi_secret
                    return fabric_ca_client
                        .register({
                            enrollmentID: 'Unionbank',
                            affiliation: 'org1.department1',
                            role: 'bank',
                            attrs: attributes
                        }, admin_user)
                })
        });

}).then((unionbank_secret) => {
    // next we need to enroll the user with CA server
    secret_unionbank = unionbank_secret
    console.log('Successfully registered Samsung - secret:' + secret_samsung);
    console.log('Successfully registered Xiaomi - secret:' + secret_xiaomi);
    console.log('Successfully registered Unionbank - secret:' + secret_unionbank);

    return fabric_ca_client
        .enroll({
            enrollmentID: 'Samsung',
            enrollmentSecret: secret_samsung
        })
        .then(() => {
            return fabric_ca_client
                .enroll({
                    enrollmentID: 'Xiaomi',
                    enrollmentSecret: secret_xiaomi
                })
                .then(() => {
                    return fabric_ca_client
                        .enroll({
                            enrollmentID: 'Unionbank',
                            enrollmentSecret: secret_unionbank
                        })
                })
        });
}).then((enrollment) => {
    console.log('Successfully enrolled member user "user1" ');
    return fabric_client
        .createUser({
            username: 'Samsung',
            mspid: 'Org1MSP',
            cryptoContent: {
                privateKeyPEM: enrollment.key.toBytes(),
                signedCertPEM: enrollment.certificate
            }
        })
        .then(() => {
            return fabric_client
                .createUser({
                    username: 'Xiaomi',
                    mspid: 'Org1MSP',
                    cryptoContent: {
                        privateKeyPEM: enrollment.key.toBytes(),
                        signedCertPEM: enrollment.certificate
                    }
                })
                .then(() => {
                    return fabric_client
                        .createUser({
                            username: 'Unionbank',
                            mspid: 'Org1MSP',
                            cryptoContent: {
                                privateKeyPEM: enrollment.key.toBytes(),
                                signedCertPEM: enrollment.certificate
                            }
                        })
                })
        });
}).then((user) => {
    member_user_samsung = user;
    member_user_xiaomi = user;
    member_user_unionbank = user;

    return fabric_client
        .setUserContext(member_user_samsung)
        .then(() => {
            return fabric_client
                .setUserContext(member_user_xiaomi)
                .then(() => {
                    return fabric_client
                        .setUserContext(member_user_unionbank)
                });
        }).then(() => {
            console.log('User1 was successfully registered and enrolled and is ready to interact with the fabric network');

        }).catch((err) => {
            console.error('Failed to register: ' + err);
            if (err.toString().indexOf('Authorization') > -1) {
                console.error('Authorization failures may be caused by having admin credentials from a previous CA instance.\n' +
                    'Try again after deleting the contents of the store directory ' + store_path);
            }
        })
});
