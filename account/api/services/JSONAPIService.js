"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
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
function error(status, detail, code, source) {
    return {
        errors: [{
                detail,
                status,
                code,
                source,
            }]
    };
}
exports.error = error;
/**
 * Returns a JSONAPI error object for a multiple errors
 *
 * @export
 * @param {Array<{status: string, detail: string, code: string, source: string}>} errors
 * @returns
 */
function errors(errors) {
    return {
        errors
    };
}
exports.errors = errors;

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbImFwaS9zZXJ2aWNlcy9KU09OQVBJU2VydmljZS50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOztBQUFBOzs7Ozs7Ozs7R0FTRztBQUNILGVBQXNCLE1BQWMsRUFBRSxNQUFjLEVBQUUsSUFBWSxFQUFFLE1BQWU7SUFDL0UsTUFBTSxDQUFDO1FBQ0gsTUFBTSxFQUFFLENBQUM7Z0JBQ0wsTUFBTTtnQkFDTixNQUFNO2dCQUNOLElBQUk7Z0JBQ0osTUFBTTthQUNULENBQUM7S0FDTCxDQUFBO0FBQ0wsQ0FBQztBQVRELHNCQVNDO0FBRUQ7Ozs7OztHQU1HO0FBQ0gsZ0JBQXVCLE1BQThFO0lBQ2pHLE1BQU0sQ0FBQztRQUNILE1BQU07S0FDVCxDQUFBO0FBQ0wsQ0FBQztBQUpELHdCQUlDIiwiZmlsZSI6ImFwaS9zZXJ2aWNlcy9KU09OQVBJU2VydmljZS5qcyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogUmV0dXJucyBhIEpTT05BUEkgZXJyb3Igb2JqZWN0IGZvciBhIHNpbmdsZSBlcnJvclxuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gc3RhdHVzIFRoZSBIVFRQIHN0YXR1cyBjb2RlXG4gKiBAcGFyYW0ge3N0cmluZ30gZGV0YWlsIFRoZSBlcnJvciBtZXNzYWdlXG4gKiBAcGFyYW0ge3N0cmluZ30gY29kZSBUaGUgdW5pcXVlIGVycm9yIGNvZGVcbiAqIEBwYXJhbSB7c3RyaW5nfSBzb3VyY2UgVGhlIGVycm9yIHNvdXJjZVxuICogQHJldHVybnMgXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBlcnJvcihzdGF0dXM6IHN0cmluZywgZGV0YWlsOiBzdHJpbmcsIGNvZGU6IHN0cmluZywgc291cmNlPzogc3RyaW5nKSB7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgZXJyb3JzOiBbe1xuICAgICAgICAgICAgZGV0YWlsLFxuICAgICAgICAgICAgc3RhdHVzLFxuICAgICAgICAgICAgY29kZSxcbiAgICAgICAgICAgIHNvdXJjZSxcbiAgICAgICAgfV1cbiAgICB9XG59XG5cbi8qKlxuICogUmV0dXJucyBhIEpTT05BUEkgZXJyb3Igb2JqZWN0IGZvciBhIG11bHRpcGxlIGVycm9yc1xuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge0FycmF5PHtzdGF0dXM6IHN0cmluZywgZGV0YWlsOiBzdHJpbmcsIGNvZGU6IHN0cmluZywgc291cmNlOiBzdHJpbmd9Pn0gZXJyb3JzIFxuICogQHJldHVybnMgXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBlcnJvcnMoZXJyb3JzOiBBcnJheTx7c3RhdHVzOiBzdHJpbmcsIGRldGFpbDogc3RyaW5nLCBjb2RlOiBzdHJpbmcsIHNvdXJjZT86IHN0cmluZ30+KSB7XG4gICAgcmV0dXJuIHtcbiAgICAgICAgZXJyb3JzXG4gICAgfVxufVxuXG4iXX0=
