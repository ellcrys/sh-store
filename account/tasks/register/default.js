/**
 * `default`
 *
 * ---------------------------------------------------------------
 *
 * This is the default Grunt tasklist that will be executed if you
 * run `grunt` in the top level directory of your app.  It is also
 * called automatically when you start Sails in development mode using
 * `sails lift` or `node app`.
 *
 * Note that when lifting your app with a custom environment setting
 * (i.e. `sails.config.environment`), Sails will look for a tasklist file
 * with the same name and run that instead of this one.
 *
 * > Note that as a special case for compatibility/historial reasons, if
 * > your environment is "production", and Sails cannot find a tasklist named
 * > `production.js`, it will attempt to run the `prod.js` tasklist as well
 * > before defaulting to `default.js`.
 *
 * For more information see:
 *   http://sailsjs.org/documentation/anatomy/my-app/tasks/register/default-js
 *
 */
module.exports = function (grunt) {
    grunt.registerTask('default', ['compileAssets', 'linkAssets', 'watch']);
};
//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInRhc2tzL3JlZ2lzdGVyL2RlZmF1bHQuanMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUE7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7R0FzQkc7QUFDSCxNQUFNLENBQUMsT0FBTyxHQUFHLFVBQVUsS0FBSztJQUM5QixLQUFLLENBQUMsWUFBWSxDQUFDLFNBQVMsRUFBRSxDQUFDLGVBQWUsRUFBRSxZQUFZLEVBQUcsT0FBTyxDQUFDLENBQUMsQ0FBQztBQUMzRSxDQUFDLENBQUMiLCJmaWxlIjoidGFza3MvcmVnaXN0ZXIvZGVmYXVsdC5qcyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogYGRlZmF1bHRgXG4gKlxuICogLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tXG4gKlxuICogVGhpcyBpcyB0aGUgZGVmYXVsdCBHcnVudCB0YXNrbGlzdCB0aGF0IHdpbGwgYmUgZXhlY3V0ZWQgaWYgeW91XG4gKiBydW4gYGdydW50YCBpbiB0aGUgdG9wIGxldmVsIGRpcmVjdG9yeSBvZiB5b3VyIGFwcC4gIEl0IGlzIGFsc29cbiAqIGNhbGxlZCBhdXRvbWF0aWNhbGx5IHdoZW4geW91IHN0YXJ0IFNhaWxzIGluIGRldmVsb3BtZW50IG1vZGUgdXNpbmdcbiAqIGBzYWlscyBsaWZ0YCBvciBgbm9kZSBhcHBgLlxuICpcbiAqIE5vdGUgdGhhdCB3aGVuIGxpZnRpbmcgeW91ciBhcHAgd2l0aCBhIGN1c3RvbSBlbnZpcm9ubWVudCBzZXR0aW5nXG4gKiAoaS5lLiBgc2FpbHMuY29uZmlnLmVudmlyb25tZW50YCksIFNhaWxzIHdpbGwgbG9vayBmb3IgYSB0YXNrbGlzdCBmaWxlXG4gKiB3aXRoIHRoZSBzYW1lIG5hbWUgYW5kIHJ1biB0aGF0IGluc3RlYWQgb2YgdGhpcyBvbmUuXG4gKiBcbiAqID4gTm90ZSB0aGF0IGFzIGEgc3BlY2lhbCBjYXNlIGZvciBjb21wYXRpYmlsaXR5L2hpc3RvcmlhbCByZWFzb25zLCBpZlxuICogPiB5b3VyIGVudmlyb25tZW50IGlzIFwicHJvZHVjdGlvblwiLCBhbmQgU2FpbHMgY2Fubm90IGZpbmQgYSB0YXNrbGlzdCBuYW1lZFxuICogPiBgcHJvZHVjdGlvbi5qc2AsIGl0IHdpbGwgYXR0ZW1wdCB0byBydW4gdGhlIGBwcm9kLmpzYCB0YXNrbGlzdCBhcyB3ZWxsXG4gKiA+IGJlZm9yZSBkZWZhdWx0aW5nIHRvIGBkZWZhdWx0LmpzYC5cbiAqXG4gKiBGb3IgbW9yZSBpbmZvcm1hdGlvbiBzZWU6XG4gKiAgIGh0dHA6Ly9zYWlsc2pzLm9yZy9kb2N1bWVudGF0aW9uL2FuYXRvbXkvbXktYXBwL3Rhc2tzL3JlZ2lzdGVyL2RlZmF1bHQtanNcbiAqXG4gKi9cbm1vZHVsZS5leHBvcnRzID0gZnVuY3Rpb24gKGdydW50KSB7XG4gIGdydW50LnJlZ2lzdGVyVGFzaygnZGVmYXVsdCcsIFsnY29tcGlsZUFzc2V0cycsICdsaW5rQXNzZXRzJywgICd3YXRjaCddKTtcbn07XG4iXX0=
//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInRhc2tzL3JlZ2lzdGVyL2RlZmF1bHQuanMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUE7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7R0FzQkc7QUFDSCxNQUFNLENBQUMsT0FBTyxHQUFHLFVBQVUsS0FBSztJQUM1QixLQUFLLENBQUMsWUFBWSxDQUFDLFNBQVMsRUFBRSxDQUFDLGVBQWUsRUFBRSxZQUFZLEVBQUUsT0FBTyxDQUFDLENBQUMsQ0FBQztBQUM1RSxDQUFDLENBQUM7QUFFRixtMURBQW0xRCIsImZpbGUiOiJ0YXNrcy9yZWdpc3Rlci9kZWZhdWx0LmpzIiwic291cmNlc0NvbnRlbnQiOlsiLyoqXG4gKiBgZGVmYXVsdGBcbiAqXG4gKiAtLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS1cbiAqXG4gKiBUaGlzIGlzIHRoZSBkZWZhdWx0IEdydW50IHRhc2tsaXN0IHRoYXQgd2lsbCBiZSBleGVjdXRlZCBpZiB5b3VcbiAqIHJ1biBgZ3J1bnRgIGluIHRoZSB0b3AgbGV2ZWwgZGlyZWN0b3J5IG9mIHlvdXIgYXBwLiAgSXQgaXMgYWxzb1xuICogY2FsbGVkIGF1dG9tYXRpY2FsbHkgd2hlbiB5b3Ugc3RhcnQgU2FpbHMgaW4gZGV2ZWxvcG1lbnQgbW9kZSB1c2luZ1xuICogYHNhaWxzIGxpZnRgIG9yIGBub2RlIGFwcGAuXG4gKlxuICogTm90ZSB0aGF0IHdoZW4gbGlmdGluZyB5b3VyIGFwcCB3aXRoIGEgY3VzdG9tIGVudmlyb25tZW50IHNldHRpbmdcbiAqIChpLmUuIGBzYWlscy5jb25maWcuZW52aXJvbm1lbnRgKSwgU2FpbHMgd2lsbCBsb29rIGZvciBhIHRhc2tsaXN0IGZpbGVcbiAqIHdpdGggdGhlIHNhbWUgbmFtZSBhbmQgcnVuIHRoYXQgaW5zdGVhZCBvZiB0aGlzIG9uZS5cbiAqXG4gKiA+IE5vdGUgdGhhdCBhcyBhIHNwZWNpYWwgY2FzZSBmb3IgY29tcGF0aWJpbGl0eS9oaXN0b3JpYWwgcmVhc29ucywgaWZcbiAqID4geW91ciBlbnZpcm9ubWVudCBpcyBcInByb2R1Y3Rpb25cIiwgYW5kIFNhaWxzIGNhbm5vdCBmaW5kIGEgdGFza2xpc3QgbmFtZWRcbiAqID4gYHByb2R1Y3Rpb24uanNgLCBpdCB3aWxsIGF0dGVtcHQgdG8gcnVuIHRoZSBgcHJvZC5qc2AgdGFza2xpc3QgYXMgd2VsbFxuICogPiBiZWZvcmUgZGVmYXVsdGluZyB0byBgZGVmYXVsdC5qc2AuXG4gKlxuICogRm9yIG1vcmUgaW5mb3JtYXRpb24gc2VlOlxuICogICBodHRwOi8vc2FpbHNqcy5vcmcvZG9jdW1lbnRhdGlvbi9hbmF0b215L215LWFwcC90YXNrcy9yZWdpc3Rlci9kZWZhdWx0LWpzXG4gKlxuICovXG5tb2R1bGUuZXhwb3J0cyA9IGZ1bmN0aW9uIChncnVudCkge1xuICAgIGdydW50LnJlZ2lzdGVyVGFzaygnZGVmYXVsdCcsIFsnY29tcGlsZUFzc2V0cycsICdsaW5rQXNzZXRzJywgJ3dhdGNoJ10pO1xufTtcblxuLy8jIHNvdXJjZU1hcHBpbmdVUkw9ZGF0YTphcHBsaWNhdGlvbi9qc29uO2NoYXJzZXQ9dXRmODtiYXNlNjQsZXlKMlpYSnphVzl1SWpvekxDSnpiM1Z5WTJWeklqcGJJblJoYzJ0ekwzSmxaMmx6ZEdWeUwyUmxabUYxYkhRdWFuTWlYU3dpYm1GdFpYTWlPbHRkTENKdFlYQndhVzVuY3lJNklrRkJRVUU3T3pzN096czdPenM3T3pzN096czdPenM3T3pzN1IwRnpRa2M3UVVGRFNDeE5RVUZOTEVOQlFVTXNUMEZCVHl4SFFVRkhMRlZCUVZVc1MwRkJTenRKUVVNNVFpeExRVUZMTEVOQlFVTXNXVUZCV1N4RFFVRkRMRk5CUVZNc1JVRkJSU3hEUVVGRExHVkJRV1VzUlVGQlJTeFpRVUZaTEVWQlFVY3NUMEZCVHl4RFFVRkRMRU5CUVVNc1EwRkJRenRCUVVNelJTeERRVUZETEVOQlFVTWlMQ0ptYVd4bElqb2lkR0Z6YTNNdmNtVm5hWE4wWlhJdlpHVm1ZWFZzZEM1cWN5SXNJbk52ZFhKalpYTkRiMjUwWlc1MElqcGJJaThxS2x4dUlDb2dZR1JsWm1GMWJIUmdYRzRnS2x4dUlDb2dMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRYRzRnS2x4dUlDb2dWR2hwY3lCcGN5QjBhR1VnWkdWbVlYVnNkQ0JIY25WdWRDQjBZWE5yYkdsemRDQjBhR0YwSUhkcGJHd2dZbVVnWlhobFkzVjBaV1FnYVdZZ2VXOTFYRzRnS2lCeWRXNGdZR2R5ZFc1MFlDQnBiaUIwYUdVZ2RHOXdJR3hsZG1Wc0lHUnBjbVZqZEc5eWVTQnZaaUI1YjNWeUlHRndjQzRnSUVsMElHbHpJR0ZzYzI5Y2JpQXFJR05oYkd4bFpDQmhkWFJ2YldGMGFXTmhiR3g1SUhkb1pXNGdlVzkxSUhOMFlYSjBJRk5oYVd4eklHbHVJR1JsZG1Wc2IzQnRaVzUwSUcxdlpHVWdkWE5wYm1kY2JpQXFJR0J6WVdsc2N5QnNhV1owWUNCdmNpQmdibTlrWlNCaGNIQmdMbHh1SUNwY2JpQXFJRTV2ZEdVZ2RHaGhkQ0IzYUdWdUlHeHBablJwYm1jZ2VXOTFjaUJoY0hBZ2QybDBhQ0JoSUdOMWMzUnZiU0JsYm5acGNtOXViV1Z1ZENCelpYUjBhVzVuWEc0Z0tpQW9hUzVsTGlCZ2MyRnBiSE11WTI5dVptbG5MbVZ1ZG1seWIyNXRaVzUwWUNrc0lGTmhhV3h6SUhkcGJHd2diRzl2YXlCbWIzSWdZU0IwWVhOcmJHbHpkQ0JtYVd4bFhHNGdLaUIzYVhSb0lIUm9aU0J6WVcxbElHNWhiV1VnWVc1a0lISjFiaUIwYUdGMElHbHVjM1JsWVdRZ2IyWWdkR2hwY3lCdmJtVXVYRzRnS2lCY2JpQXFJRDRnVG05MFpTQjBhR0YwSUdGeklHRWdjM0JsWTJsaGJDQmpZWE5sSUdadmNpQmpiMjF3WVhScFltbHNhWFI1TDJocGMzUnZjbWxoYkNCeVpXRnpiMjV6TENCcFpseHVJQ29nUGlCNWIzVnlJR1Z1ZG1seWIyNXRaVzUwSUdseklGd2ljSEp2WkhWamRHbHZibHdpTENCaGJtUWdVMkZwYkhNZ1kyRnVibTkwSUdacGJtUWdZU0IwWVhOcmJHbHpkQ0J1WVcxbFpGeHVJQ29nUGlCZ2NISnZaSFZqZEdsdmJpNXFjMkFzSUdsMElIZHBiR3dnWVhSMFpXMXdkQ0IwYnlCeWRXNGdkR2hsSUdCd2NtOWtMbXB6WUNCMFlYTnJiR2x6ZENCaGN5QjNaV3hzWEc0Z0tpQStJR0psWm05eVpTQmtaV1poZFd4MGFXNW5JSFJ2SUdCa1pXWmhkV3gwTG1wellDNWNiaUFxWEc0Z0tpQkdiM0lnYlc5eVpTQnBibVp2Y20xaGRHbHZiaUJ6WldVNlhHNGdLaUFnSUdoMGRIQTZMeTl6WVdsc2MycHpMbTl5Wnk5a2IyTjFiV1Z1ZEdGMGFXOXVMMkZ1WVhSdmJYa3ZiWGt0WVhCd0wzUmhjMnR6TDNKbFoybHpkR1Z5TDJSbFptRjFiSFF0YW5OY2JpQXFYRzRnS2k5Y2JtMXZaSFZzWlM1bGVIQnZjblJ6SUQwZ1puVnVZM1JwYjI0Z0tHZHlkVzUwS1NCN1hHNGdJR2R5ZFc1MExuSmxaMmx6ZEdWeVZHRnpheWduWkdWbVlYVnNkQ2NzSUZzblkyOXRjR2xzWlVGemMyVjBjeWNzSUNkc2FXNXJRWE56WlhSekp5d2dJQ2QzWVhSamFDZGRLVHRjYm4wN1hHNGlYWDA9XG4iXX0=

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInRhc2tzL3JlZ2lzdGVyL2RlZmF1bHQuanMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUE7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7R0FzQkc7QUFDSCxNQUFNLENBQUMsT0FBTyxHQUFHLFVBQVUsS0FBSztJQUM1QixLQUFLLENBQUMsWUFBWSxDQUFDLFNBQVMsRUFBRSxDQUFDLGVBQWUsRUFBRSxZQUFZLEVBQUUsT0FBTyxDQUFDLENBQUMsQ0FBQztBQUM1RSxDQUFDLENBQUM7QUFDRixtMURBQW0xRDtBQUVuMUQsK3lJQUEreUkiLCJmaWxlIjoidGFza3MvcmVnaXN0ZXIvZGVmYXVsdC5qcyIsInNvdXJjZXNDb250ZW50IjpbIi8qKlxuICogYGRlZmF1bHRgXG4gKlxuICogLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tXG4gKlxuICogVGhpcyBpcyB0aGUgZGVmYXVsdCBHcnVudCB0YXNrbGlzdCB0aGF0IHdpbGwgYmUgZXhlY3V0ZWQgaWYgeW91XG4gKiBydW4gYGdydW50YCBpbiB0aGUgdG9wIGxldmVsIGRpcmVjdG9yeSBvZiB5b3VyIGFwcC4gIEl0IGlzIGFsc29cbiAqIGNhbGxlZCBhdXRvbWF0aWNhbGx5IHdoZW4geW91IHN0YXJ0IFNhaWxzIGluIGRldmVsb3BtZW50IG1vZGUgdXNpbmdcbiAqIGBzYWlscyBsaWZ0YCBvciBgbm9kZSBhcHBgLlxuICpcbiAqIE5vdGUgdGhhdCB3aGVuIGxpZnRpbmcgeW91ciBhcHAgd2l0aCBhIGN1c3RvbSBlbnZpcm9ubWVudCBzZXR0aW5nXG4gKiAoaS5lLiBgc2FpbHMuY29uZmlnLmVudmlyb25tZW50YCksIFNhaWxzIHdpbGwgbG9vayBmb3IgYSB0YXNrbGlzdCBmaWxlXG4gKiB3aXRoIHRoZSBzYW1lIG5hbWUgYW5kIHJ1biB0aGF0IGluc3RlYWQgb2YgdGhpcyBvbmUuXG4gKlxuICogPiBOb3RlIHRoYXQgYXMgYSBzcGVjaWFsIGNhc2UgZm9yIGNvbXBhdGliaWxpdHkvaGlzdG9yaWFsIHJlYXNvbnMsIGlmXG4gKiA+IHlvdXIgZW52aXJvbm1lbnQgaXMgXCJwcm9kdWN0aW9uXCIsIGFuZCBTYWlscyBjYW5ub3QgZmluZCBhIHRhc2tsaXN0IG5hbWVkXG4gKiA+IGBwcm9kdWN0aW9uLmpzYCwgaXQgd2lsbCBhdHRlbXB0IHRvIHJ1biB0aGUgYHByb2QuanNgIHRhc2tsaXN0IGFzIHdlbGxcbiAqID4gYmVmb3JlIGRlZmF1bHRpbmcgdG8gYGRlZmF1bHQuanNgLlxuICpcbiAqIEZvciBtb3JlIGluZm9ybWF0aW9uIHNlZTpcbiAqICAgaHR0cDovL3NhaWxzanMub3JnL2RvY3VtZW50YXRpb24vYW5hdG9teS9teS1hcHAvdGFza3MvcmVnaXN0ZXIvZGVmYXVsdC1qc1xuICpcbiAqL1xubW9kdWxlLmV4cG9ydHMgPSBmdW5jdGlvbiAoZ3J1bnQpIHtcbiAgICBncnVudC5yZWdpc3RlclRhc2soJ2RlZmF1bHQnLCBbJ2NvbXBpbGVBc3NldHMnLCAnbGlua0Fzc2V0cycsICd3YXRjaCddKTtcbn07XG4vLyMgc291cmNlTWFwcGluZ1VSTD1kYXRhOmFwcGxpY2F0aW9uL2pzb247Y2hhcnNldD11dGY4O2Jhc2U2NCxleUoyWlhKemFXOXVJam96TENKemIzVnlZMlZ6SWpwYkluUmhjMnR6TDNKbFoybHpkR1Z5TDJSbFptRjFiSFF1YW5NaVhTd2libUZ0WlhNaU9sdGRMQ0p0WVhCd2FXNW5jeUk2SWtGQlFVRTdPenM3T3pzN096czdPenM3T3pzN096czdPenM3UjBGelFrYzdRVUZEU0N4TlFVRk5MRU5CUVVNc1QwRkJUeXhIUVVGSExGVkJRVlVzUzBGQlN6dEpRVU01UWl4TFFVRkxMRU5CUVVNc1dVRkJXU3hEUVVGRExGTkJRVk1zUlVGQlJTeERRVUZETEdWQlFXVXNSVUZCUlN4WlFVRlpMRVZCUVVjc1QwRkJUeXhEUVVGRExFTkJRVU1zUTBGQlF6dEJRVU16UlN4RFFVRkRMRU5CUVVNaUxDSm1hV3hsSWpvaWRHRnphM012Y21WbmFYTjBaWEl2WkdWbVlYVnNkQzVxY3lJc0luTnZkWEpqWlhORGIyNTBaVzUwSWpwYklpOHFLbHh1SUNvZ1lHUmxabUYxYkhSZ1hHNGdLbHh1SUNvZ0xTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdFhHNGdLbHh1SUNvZ1ZHaHBjeUJwY3lCMGFHVWdaR1ZtWVhWc2RDQkhjblZ1ZENCMFlYTnJiR2x6ZENCMGFHRjBJSGRwYkd3Z1ltVWdaWGhsWTNWMFpXUWdhV1lnZVc5MVhHNGdLaUJ5ZFc0Z1lHZHlkVzUwWUNCcGJpQjBhR1VnZEc5d0lHeGxkbVZzSUdScGNtVmpkRzl5ZVNCdlppQjViM1Z5SUdGd2NDNGdJRWwwSUdseklHRnNjMjljYmlBcUlHTmhiR3hsWkNCaGRYUnZiV0YwYVdOaGJHeDVJSGRvWlc0Z2VXOTFJSE4wWVhKMElGTmhhV3h6SUdsdUlHUmxkbVZzYjNCdFpXNTBJRzF2WkdVZ2RYTnBibWRjYmlBcUlHQnpZV2xzY3lCc2FXWjBZQ0J2Y2lCZ2JtOWtaU0JoY0hCZ0xseHVJQ3BjYmlBcUlFNXZkR1VnZEdoaGRDQjNhR1Z1SUd4cFpuUnBibWNnZVc5MWNpQmhjSEFnZDJsMGFDQmhJR04xYzNSdmJTQmxiblpwY205dWJXVnVkQ0J6WlhSMGFXNW5YRzRnS2lBb2FTNWxMaUJnYzJGcGJITXVZMjl1Wm1sbkxtVnVkbWx5YjI1dFpXNTBZQ2tzSUZOaGFXeHpJSGRwYkd3Z2JHOXZheUJtYjNJZ1lTQjBZWE5yYkdsemRDQm1hV3hsWEc0Z0tpQjNhWFJvSUhSb1pTQnpZVzFsSUc1aGJXVWdZVzVrSUhKMWJpQjBhR0YwSUdsdWMzUmxZV1FnYjJZZ2RHaHBjeUJ2Ym1VdVhHNGdLaUJjYmlBcUlENGdUbTkwWlNCMGFHRjBJR0Z6SUdFZ2MzQmxZMmxoYkNCallYTmxJR1p2Y2lCamIyMXdZWFJwWW1sc2FYUjVMMmhwYzNSdmNtbGhiQ0J5WldGemIyNXpMQ0JwWmx4dUlDb2dQaUI1YjNWeUlHVnVkbWx5YjI1dFpXNTBJR2x6SUZ3aWNISnZaSFZqZEdsdmJsd2lMQ0JoYm1RZ1UyRnBiSE1nWTJGdWJtOTBJR1pwYm1RZ1lTQjBZWE5yYkdsemRDQnVZVzFsWkZ4dUlDb2dQaUJnY0hKdlpIVmpkR2x2Ymk1cWMyQXNJR2wwSUhkcGJHd2dZWFIwWlcxd2RDQjBieUJ5ZFc0Z2RHaGxJR0J3Y205a0xtcHpZQ0IwWVhOcmJHbHpkQ0JoY3lCM1pXeHNYRzRnS2lBK0lHSmxabTl5WlNCa1pXWmhkV3gwYVc1bklIUnZJR0JrWldaaGRXeDBMbXB6WUM1Y2JpQXFYRzRnS2lCR2IzSWdiVzl5WlNCcGJtWnZjbTFoZEdsdmJpQnpaV1U2WEc0Z0tpQWdJR2gwZEhBNkx5OXpZV2xzYzJwekxtOXlaeTlrYjJOMWJXVnVkR0YwYVc5dUwyRnVZWFJ2YlhrdmJYa3RZWEJ3TDNSaGMydHpMM0psWjJsemRHVnlMMlJsWm1GMWJIUXRhbk5jYmlBcVhHNGdLaTljYm0xdlpIVnNaUzVsZUhCdmNuUnpJRDBnWm5WdVkzUnBiMjRnS0dkeWRXNTBLU0I3WEc0Z0lHZHlkVzUwTG5KbFoybHpkR1Z5VkdGemF5Z25aR1ZtWVhWc2RDY3NJRnNuWTI5dGNHbHNaVUZ6YzJWMGN5Y3NJQ2RzYVc1clFYTnpaWFJ6Snl3Z0lDZDNZWFJqYUNkZEtUdGNibjA3WEc0aVhYMD1cblxuLy8jIHNvdXJjZU1hcHBpbmdVUkw9ZGF0YTphcHBsaWNhdGlvbi9qc29uO2NoYXJzZXQ9dXRmODtiYXNlNjQsZXlKMlpYSnphVzl1SWpvekxDSnpiM1Z5WTJWeklqcGJJblJoYzJ0ekwzSmxaMmx6ZEdWeUwyUmxabUYxYkhRdWFuTWlYU3dpYm1GdFpYTWlPbHRkTENKdFlYQndhVzVuY3lJNklrRkJRVUU3T3pzN096czdPenM3T3pzN096czdPenM3T3pzN1IwRnpRa2M3UVVGRFNDeE5RVUZOTEVOQlFVTXNUMEZCVHl4SFFVRkhMRlZCUVZVc1MwRkJTenRKUVVNMVFpeExRVUZMTEVOQlFVTXNXVUZCV1N4RFFVRkRMRk5CUVZNc1JVRkJSU3hEUVVGRExHVkJRV1VzUlVGQlJTeFpRVUZaTEVWQlFVVXNUMEZCVHl4RFFVRkRMRU5CUVVNc1EwRkJRenRCUVVNMVJTeERRVUZETEVOQlFVTTdRVUZGUml4dE1VUkJRVzB4UkNJc0ltWnBiR1VpT2lKMFlYTnJjeTl5WldkcGMzUmxjaTlrWldaaGRXeDBMbXB6SWl3aWMyOTFjbU5sYzBOdmJuUmxiblFpT2xzaUx5b3FYRzRnS2lCZ1pHVm1ZWFZzZEdCY2JpQXFYRzRnS2lBdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMwdExTMHRMUzB0TFMxY2JpQXFYRzRnS2lCVWFHbHpJR2x6SUhSb1pTQmtaV1poZFd4MElFZHlkVzUwSUhSaGMydHNhWE4wSUhSb1lYUWdkMmxzYkNCaVpTQmxlR1ZqZFhSbFpDQnBaaUI1YjNWY2JpQXFJSEoxYmlCZ1ozSjFiblJnSUdsdUlIUm9aU0IwYjNBZ2JHVjJaV3dnWkdseVpXTjBiM0o1SUc5bUlIbHZkWElnWVhCd0xpQWdTWFFnYVhNZ1lXeHpiMXh1SUNvZ1kyRnNiR1ZrSUdGMWRHOXRZWFJwWTJGc2JIa2dkMmhsYmlCNWIzVWdjM1JoY25RZ1UyRnBiSE1nYVc0Z1pHVjJaV3h2Y0cxbGJuUWdiVzlrWlNCMWMybHVaMXh1SUNvZ1lITmhhV3h6SUd4cFpuUmdJRzl5SUdCdWIyUmxJR0Z3Y0dBdVhHNGdLbHh1SUNvZ1RtOTBaU0IwYUdGMElIZG9aVzRnYkdsbWRHbHVaeUI1YjNWeUlHRndjQ0IzYVhSb0lHRWdZM1Z6ZEc5dElHVnVkbWx5YjI1dFpXNTBJSE5sZEhScGJtZGNiaUFxSUNocExtVXVJR0J6WVdsc2N5NWpiMjVtYVdjdVpXNTJhWEp2Ym0xbGJuUmdLU3dnVTJGcGJITWdkMmxzYkNCc2IyOXJJR1p2Y2lCaElIUmhjMnRzYVhOMElHWnBiR1ZjYmlBcUlIZHBkR2dnZEdobElITmhiV1VnYm1GdFpTQmhibVFnY25WdUlIUm9ZWFFnYVc1emRHVmhaQ0J2WmlCMGFHbHpJRzl1WlM1Y2JpQXFYRzRnS2lBK0lFNXZkR1VnZEdoaGRDQmhjeUJoSUhOd1pXTnBZV3dnWTJGelpTQm1iM0lnWTI5dGNHRjBhV0pwYkdsMGVTOW9hWE4wYjNKcFlXd2djbVZoYzI5dWN5d2dhV1pjYmlBcUlENGdlVzkxY2lCbGJuWnBjbTl1YldWdWRDQnBjeUJjSW5CeWIyUjFZM1JwYjI1Y0lpd2dZVzVrSUZOaGFXeHpJR05oYm01dmRDQm1hVzVrSUdFZ2RHRnphMnhwYzNRZ2JtRnRaV1JjYmlBcUlENGdZSEJ5YjJSMVkzUnBiMjR1YW5OZ0xDQnBkQ0IzYVd4c0lHRjBkR1Z0Y0hRZ2RHOGdjblZ1SUhSb1pTQmdjSEp2WkM1cWMyQWdkR0Z6YTJ4cGMzUWdZWE1nZDJWc2JGeHVJQ29nUGlCaVpXWnZjbVVnWkdWbVlYVnNkR2x1WnlCMGJ5QmdaR1ZtWVhWc2RDNXFjMkF1WEc0Z0tseHVJQ29nUm05eUlHMXZjbVVnYVc1bWIzSnRZWFJwYjI0Z2MyVmxPbHh1SUNvZ0lDQm9kSFJ3T2k4dmMyRnBiSE5xY3k1dmNtY3ZaRzlqZFcxbGJuUmhkR2x2Ymk5aGJtRjBiMjE1TDIxNUxXRndjQzkwWVhOcmN5OXlaV2RwYzNSbGNpOWtaV1poZFd4MExXcHpYRzRnS2x4dUlDb3ZYRzV0YjJSMWJHVXVaWGh3YjNKMGN5QTlJR1oxYm1OMGFXOXVJQ2huY25WdWRDa2dlMXh1SUNBZ0lHZHlkVzUwTG5KbFoybHpkR1Z5VkdGemF5Z25aR1ZtWVhWc2RDY3NJRnNuWTI5dGNHbHNaVUZ6YzJWMGN5Y3NJQ2RzYVc1clFYTnpaWFJ6Snl3Z0ozZGhkR05vSjEwcE8xeHVmVHRjYmx4dUx5OGpJSE52ZFhKalpVMWhjSEJwYm1kVlVrdzlaR0YwWVRwaGNIQnNhV05oZEdsdmJpOXFjMjl1TzJOb1lYSnpaWFE5ZFhSbU9EdGlZWE5sTmpRc1pYbEtNbHBZU25waFZ6bDFTV3B2ZWt4RFNucGlNMVo1V1RKV2VrbHFjR0pKYmxKb1l6SjBla3d6U214YU1teDZaRWRXZVV3eVVteGFiVVl4WWtoUmRXRnVUV2xZVTNkcFltMUdkRnBZVFdsUGJIUmtURU5LZEZsWVFuZGhWelZ1WTNsSk5rbHJSa0pSVlVVM1QzcHpOMDk2Y3pkUGVuTTNUM3B6TjA5NmN6ZFBlbk0zVDNwek4xSXdSbnBSYTJNM1VWVkdSRk5EZUU1UlZVWk9URVZPUWxGVlRYTlVNRVpDVkhsNFNGRlZSa2hNUmxaQ1VWWlZjMU13UmtKVGVuUktVVlZOTlZGcGVFeFJWVVpNVEVWT1FsRlZUWE5YVlVaQ1YxTjRSRkZWUmtSTVJrNUNVVlpOYzFKVlJrSlNVM2hFVVZWR1JFeEhWa0pSVjFWelVsVkdRbEpUZUZwUlZVWmFURVZXUWxGVlkzTlVNRVpDVkhsNFJGRlZSa1JNUlU1Q1VWVk5jMUV3UmtKUmVuUkNVVlZOZWxKVGVFUlJWVVpFVEVWT1FsRlZUV2xNUTBwdFlWZDRiRWxxYjJsa1IwWjZZVE5OZG1OdFZtNWhXRTR3V2xoSmRscEhWbTFaV0ZaelpFTTFjV041U1hOSmJrNTJaRmhLYWxwWVRrUmlNalV3V2xjMU1FbHFjR0pKYVRoeFMyeDRkVWxEYjJkWlIxSnNXbTFHTVdKSVVtZFlSelJuUzJ4NGRVbERiMmRNVXpCMFRGTXdkRXhUTUhSTVV6QjBURk13ZEV4VE1IUk1VekIwVEZNd2RFeFRNSFJNVXpCMFRGTXdkRXhUTUhSTVV6QjBURk13ZEV4VE1IUk1VekIwVEZNd2RFeFRNSFJNVXpCMFRGTXdkRXhUTUhSWVJ6Um5TMng0ZFVsRGIyZFdSMmh3WTNsQ2NHTjVRakJoUjFWbldrZFdiVmxZVm5Oa1EwSklZMjVXZFdSRFFqQlpXRTV5WWtkc2VtUkRRakJoUjBZd1NVaGtjR0pIZDJkWmJWVm5XbGhvYkZrelZqQmFWMUZuWVZkWloyVlhPVEZZUnpSblMybENlV1JYTkdkWlIyUjVaRmMxTUZsRFFuQmlhVUl3WVVkVloyUkhPWGRKUjNoc1pHMVdjMGxIVW5CamJWWnFaRWM1ZVdWVFFuWmFhVUkxWWpOV2VVbEhSbmRqUXpSblNVVnNNRWxIYkhwSlIwWnpZekk1WTJKcFFYRkpSMDVvWWtkNGJGcERRbWhrV0ZKMllsZEdNR0ZYVG1oaVIzZzFTVWhrYjFwWE5HZGxWemt4U1VoT01GbFlTakJKUms1b1lWZDRla2xIYkhWSlIxSnNaRzFXYzJJelFuUmFWelV3U1VjeGRscEhWV2RrV0U1d1ltMWtZMkpwUVhGSlIwSjZXVmRzYzJONVFuTmhWMW93V1VOQ2RtTnBRbWRpYlRscldsTkNhR05JUW1kTWJIaDFTVU53WTJKcFFYRkpSVFYyWkVkVloyUkhhR2hrUTBJellVZFdkVWxIZUhCYWJsSndZbTFqWjJWWE9URmphVUpvWTBoQloyUXliREJoUTBKb1NVZE9NV016VW5aaVUwSnNZbTVhY0dOdE9YVmlWMVoxWkVOQ2VscFlVakJoVnpWdVdFYzBaMHRwUVc5aFV6VnNUR2xDWjJNeVJuQmlTRTExV1RJNWRWcHRiRzVNYlZaMVpHMXNlV0l5TlhSYVZ6VXdXVU5yYzBsR1RtaGhWM2g2U1Voa2NHSkhkMmRpUnpsMllYbENiV0l6U1dkWlUwSXdXVmhPY21KSGJIcGtRMEp0WVZkNGJGaEhOR2RMYVVJellWaFNiMGxJVW05YVUwSjZXVmN4YkVsSE5XaGlWMVZuV1ZjMWEwbElTakZpYVVJd1lVZEdNRWxIYkhWak0xSnNXVmRSWjJJeVdXZGtSMmh3WTNsQ2RtSnRWWFZZUnpSblMybENZMkpwUVhGSlJEUm5WRzA1TUZwVFFqQmhSMFl3U1VkR2VrbEhSV2RqTTBKc1dUSnNhR0pEUW1wWldFNXNTVWRhZG1OcFFtcGlNakYzV1ZoU2NGbHRiSE5oV0ZJMVRESm9jR016VW5aamJXeG9Za05DZVZwWFJucGlNalY2VEVOQ2NGcHNlSFZKUTI5blVHbENOV0l6Vm5sSlIxWjFaRzFzZVdJeU5YUmFWelV3U1Vkc2VrbEdkMmxqU0VwMldraFdhbVJIYkhaaWJIZHBURU5DYUdKdFVXZFZNa1p3WWtoTloxa3lSblZpYlRrd1NVZGFjR0p0VVdkWlUwSXdXVmhPY21KSGJIcGtRMEoxV1ZjeGJGcEdlSFZKUTI5blVHbENaMk5JU25aYVNGWnFaRWRzZG1KcE5YRmpNa0Z6U1Vkc01FbElaSEJpUjNkbldWaFNNRnBYTVhka1EwSXdZbmxDZVdSWE5HZGtSMmhzU1VkQ2QyTnRPV3RNYlhCNldVTkNNRmxZVG5KaVIyeDZaRU5DYUdONVFqTmFWM2h6V0VjMFowdHBRU3RKUjBwc1dtMDVlVnBUUW10YVYxcG9aRmQ0TUdGWE5XNUpTRkoyU1VkQ2ExcFhXbWhrVjNnd1RHMXdlbGxETldOaWFVRnhXRWMwWjB0cFFrZGlNMGxuWWxjNWVWcFRRbkJpYlZwMlkyMHhhR1JIYkhaaWFVSjZXbGRWTmxoSE5HZExhVUZuU1Vkb01HUklRVFpNZVRsNldWZHNjMk15Y0hwTWJUbDVXbms1YTJJeVRqRmlWMVoxWkVkR01HRlhPWFZNTWtaMVdWaFNkbUpZYTNaaVdHdDBXVmhDZDB3elVtaGpNblI2VEROS2JGb3liSHBrUjFaNVRESlNiRnB0UmpGaVNGRjBZVzVPWTJKcFFYRllSelJuUzJrNVkySnRNWFphU0ZaeldsTTFiR1ZJUW5aamJsSjZTVVF3WjFwdVZuVlpNMUp3WWpJMFowdEhaSGxrVnpVd1MxTkNOMWhITkdkSlIyUjVaRmMxTUV4dVNteGFNbXg2WkVkV2VWWkhSbnBoZVdkdVdrZFdiVmxZVm5Oa1EyTnpTVVp6YmxreU9YUmpSMnh6V2xWR2VtTXlWakJqZVdOelNVTmtjMkZYTlhKUldFNTZXbGhTZWtwNWQyZEpRMlF6V1ZoU2FtRkRaR1JMVkhSalltNHdOMWhITkdsWVdEQTlYRzRpWFgwPVxuIl19