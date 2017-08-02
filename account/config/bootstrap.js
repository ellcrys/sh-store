/**
* Bootstrap
* (sails.config.bootstrap)
*
* An asynchronous bootstrap function that runs before your Sails app gets lifted.
* This gives you an opportunity to set up your data model, run jobs, or perform some special logic.
*
* For more information on bootstrapping your app, check out:
* http://sailsjs.org/#!/documentation/reference/sails.config/sails.config.bootstrap.html
*/
var db = require('./db');
var flash = require('express-flash')

module.exports.bootstrap = function (cb) {
	
	sails.hooks.http.app.use(flash())
    
    // connect to database and set connection on sails global object
    sails.db = db.connectToDatabase()
    
    // set a function to return the connection on sails global object
    sails.getDB = () => {
        return sails.db
    }
    
    // define (create) tables
    db.defineModels(sails.db).then(function () {
        console.log("Connected to database. Tables created");
        cb();
    }).catch(function (e) {
        console.error("failed to connect and setup database.", e.message);
        process.exit(1);
    });
};
