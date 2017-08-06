import {castDB} from './helpers'
const Promise = require("bluebird")


/**
 * Get contract
 * 
 * @export
 * @param {string} email 
 * @returns {Promise<any>} 
 */
export function get(q: any): Promise<any> {
    return new Promise((resolve, reject) => {
        return castDB(sails.db).models.contract.findOne({where: q}).then(resolve).catch(reject)
    })
}

/**
 * Create a contract
 * 
 * @export
 * @param {Contract} obj 
 * @returns {Promise<any>} 
 */
export function create(obj: models.Contract): Promise<any> {
    return new Promise((resolve, reject) => {
        return castDB(sails.db).models.contract.create(obj).then(resolve).catch(reject)
    })
} 