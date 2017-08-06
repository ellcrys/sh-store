/**
 * Returns a JSONAPI error object for a single error
 * 
 * @export
 * @param {string} status The HTTP status code
 * @param {string} detail The error message
 * @param {string} code The unique error code
 * @param {string} source The error source
 * @returns 
 */
export function error(status: string, detail: string, code: string, source?: string) {
    return {
        errors: [{
            detail,
            status,
            code,
            source,
        }]
    }
}

/**
 * Returns a JSONAPI error object for a multiple errors
 * 
 * @export
 * @param {Array<{status: string, detail: string, code: string, source: string}>} errors 
 * @returns 
 */
export function errors(errors: Array<{status: string, detail: string, code: string, source?: string}>) {
    return {
        errors
    }
}

