"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/**
 * Returns retVal if a and b are the same
 *
 * @export
 * @param {string} a The field to check against
 * @param {string} b The error affected field
 * @param {string} retVal The value to return if a and b are equal
 * @returns {string} Returns retVal or nothing
 */
function eqlThenReturn(a, b, retVal) {
    return (a == b) ? retVal : "";
}
exports.eqlThenReturn = eqlThenReturn;
/**
 * Log a value
 * @param v
 */
function log(v) {
    console.log(v);
}
exports.log = log;

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInZpZXdzL2hlbHBlcnMvaW5kZXgudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFBQTs7Ozs7Ozs7R0FRRztBQUNILHVCQUE4QixDQUFTLEVBQUUsQ0FBUyxFQUFFLE1BQWM7SUFDOUQsTUFBTSxDQUFDLENBQUMsQ0FBQyxJQUFJLENBQUMsQ0FBQyxHQUFHLE1BQU0sR0FBRSxFQUFFLENBQUE7QUFDaEMsQ0FBQztBQUZELHNDQUVDO0FBRUQ7OztHQUdHO0FBQ0gsYUFBb0IsQ0FBTTtJQUN0QixPQUFPLENBQUMsR0FBRyxDQUFDLENBQUMsQ0FBQyxDQUFBO0FBQ2xCLENBQUM7QUFGRCxrQkFFQyIsImZpbGUiOiJ2aWV3cy9oZWxwZXJzL2luZGV4LmpzIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBSZXR1cm5zIHJldFZhbCBpZiBhIGFuZCBiIGFyZSB0aGUgc2FtZVxuICogXG4gKiBAZXhwb3J0XG4gKiBAcGFyYW0ge3N0cmluZ30gYSBUaGUgZmllbGQgdG8gY2hlY2sgYWdhaW5zdFxuICogQHBhcmFtIHtzdHJpbmd9IGIgVGhlIGVycm9yIGFmZmVjdGVkIGZpZWxkXG4gKiBAcGFyYW0ge3N0cmluZ30gcmV0VmFsIFRoZSB2YWx1ZSB0byByZXR1cm4gaWYgYSBhbmQgYiBhcmUgZXF1YWxcbiAqIEByZXR1cm5zIHtzdHJpbmd9IFJldHVybnMgcmV0VmFsIG9yIG5vdGhpbmdcbiAqL1xuZXhwb3J0IGZ1bmN0aW9uIGVxbFRoZW5SZXR1cm4oYTogc3RyaW5nLCBiOiBzdHJpbmcsIHJldFZhbDogc3RyaW5nKTogc3RyaW5nIHtcbiAgICByZXR1cm4gKGEgPT0gYikgPyByZXRWYWw6IFwiXCJcbn1cblxuLyoqXG4gKiBMb2cgYSB2YWx1ZVxuICogQHBhcmFtIHYgXG4gKi9cbmV4cG9ydCBmdW5jdGlvbiBsb2codjogYW55KSB7XG4gICAgY29uc29sZS5sb2codilcbn0iXX0=
