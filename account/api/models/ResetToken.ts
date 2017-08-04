import {castDB} from './helpers'
const Promise = require("bluebird")

/**
 * Create a reset token
 * 
 * @export
 * @param {Account} obj 
 * @returns {Promise<any>} 
 */
export function create(obj: { token: string, account: string }): Promise<any> {
    return new Promise((resolve, reject) => {
        return castDB(sails.db).models.reset_token.create(obj).then(resolve).catch(reject)
    })
} 

/**
 * Get reset token
 * 
 * @export
 * @param {string} q the query 
 * @returns {Promise<any>} 
 */
export function get(q: any): Promise<any> {
    return new Promise((resolve, reject) => {
        return castDB(sails.db).models.reset_token.findOne({where: q}).then(resolve).catch(reject)
    })
}
