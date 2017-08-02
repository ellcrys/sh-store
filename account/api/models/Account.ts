import {castDB} from './helpers'
const Promise = require("bluebird")


/**
 * Get account by email
 * 
 * @export
 * @param {string} email 
 * @returns {Promise<any>} 
 */
export function get(q: any): Promise<any> {
    return new Promise((resolve, reject) => {
        return castDB(sails.db).models.account.findOne({where: q}).then(resolve).catch(reject)
    })
}

/**
 * Create an account
 * 
 * @export
 * @param {Account} obj 
 * @returns {Promise<any>} 
 */
export function create(obj: Account): Promise<any> {
    return new Promise((resolve, reject) => {
        return castDB(sails.db).models.account.create(obj).then(resolve).catch(reject)
    })
} 