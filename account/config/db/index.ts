import Sequelize = require('sequelize')
import SequelizeCockroachDB = require('sequelize-cockroachdb');
import Promise = require('bluebird');

/**
 * Connect to database
 * @export
 * @returns {Sequelize} Sequelize instance
 */
export function connectToDatabase(): Sequelize.Sequelize {
    var options = sails.config.connections.postgres
    var sequelize = new SequelizeCockroachDB(options.database, options.user, options.password, {
        dialect: 'postgres',
        host: options.host,
        port: options.port,
        logging: options.logging
    })
    return sequelize
}

/**
 * Define models. Create tables if not existing.
 * 
 * @export
 * @param {Sequelize} sequelize 
 * @returns {Promise}
 */
export function defineModels(sequelize: Sequelize.Sequelize): Promise<any> {

    sequelize.define('account', {
        sn: { type: Sequelize.BIGINT, primaryKey: true, autoIncrement: true },
        id: { type: Sequelize.STRING },
        first_name: { type: Sequelize.STRING },
        last_name: { type: Sequelize.STRING },
        email: { type: Sequelize.STRING },
        password: { type: Sequelize.STRING },
        developer: { type: Sequelize.BOOLEAN, defaultValue: false },
        confirmed: { type: Sequelize.BOOLEAN, defaultValue: false },
        client_id: { type: Sequelize.STRING },
        client_secret: { type: Sequelize.STRING },
        confirmation_code: { type: Sequelize.STRING }
    }, { underscored: true })

    return sequelize.sync({})
}