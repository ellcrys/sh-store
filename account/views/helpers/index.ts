/**
 * Returns retVal if a and b are the same
 * 
 * @export
 * @param {string} a The field to check against
 * @param {string} b The error affected field
 * @param {string} retVal The value to return if a and b are equal
 * @returns {string} Returns retVal or nothing
 */
export function eqlThenReturn(a: string, b: string, retVal: string): string {
    return (a == b) ? retVal: ""
}

/**
 * Log a value
 * @param v 
 */
export function log(v: any) {
    console.log(v)
}