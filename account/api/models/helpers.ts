import * as sequelize from 'sequelize';

/**
 * Cast a database object to sequelize.Sequelize
 * @param {*} db 
 * @returns {sequelize.Sequelize} 
 */
export function castDB(db: any): sequelize.Sequelize {
    return (db as sequelize.Sequelize)
}