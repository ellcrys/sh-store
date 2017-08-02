define(function () {
    "use strict";
    return function (elem) {
        // Support: IE <=11 only, Firefox <=30 (#15098, #14150)
        // IE throws on elements created in popups
        // FF meanwhile throws on frame elements through "defaultView.getComputedStyle"
        var view = elem.ownerDocument.defaultView;
        if (!view || !view.opener) {
            view = window;
        }
        return view.getComputedStyle(elem);
    };
});
//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInd3dy9wYWNrcy9qcXVlcnkvc3JjL2Nzcy92YXIvZ2V0U3R5bGVzLmpzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBLE1BQU0sQ0FBRTtJQUNQLFlBQVksQ0FBQztJQUViLE1BQU0sQ0FBQyxVQUFVLElBQUk7UUFFcEIsdURBQXVEO1FBQ3ZELDBDQUEwQztRQUMxQywrRUFBK0U7UUFDL0UsSUFBSSxJQUFJLEdBQUcsSUFBSSxDQUFDLGFBQWEsQ0FBQyxXQUFXLENBQUM7UUFFMUMsRUFBRSxDQUFDLENBQUUsQ0FBQyxJQUFJLElBQUksQ0FBQyxJQUFJLENBQUMsTUFBTyxDQUFDLENBQUMsQ0FBQztZQUM3QixJQUFJLEdBQUcsTUFBTSxDQUFDO1FBQ2YsQ0FBQztRQUVELE1BQU0sQ0FBQyxJQUFJLENBQUMsZ0JBQWdCLENBQUUsSUFBSSxDQUFFLENBQUM7SUFDdEMsQ0FBQyxDQUFDO0FBQ0gsQ0FBQyxDQUFFLENBQUMiLCJmaWxlIjoid3d3L3BhY2tzL2pxdWVyeS9zcmMvY3NzL3Zhci9nZXRTdHlsZXMuanMiLCJzb3VyY2VzQ29udGVudCI6WyJkZWZpbmUoIGZ1bmN0aW9uKCkge1xuXHRcInVzZSBzdHJpY3RcIjtcblxuXHRyZXR1cm4gZnVuY3Rpb24oIGVsZW0gKSB7XG5cblx0XHQvLyBTdXBwb3J0OiBJRSA8PTExIG9ubHksIEZpcmVmb3ggPD0zMCAoIzE1MDk4LCAjMTQxNTApXG5cdFx0Ly8gSUUgdGhyb3dzIG9uIGVsZW1lbnRzIGNyZWF0ZWQgaW4gcG9wdXBzXG5cdFx0Ly8gRkYgbWVhbndoaWxlIHRocm93cyBvbiBmcmFtZSBlbGVtZW50cyB0aHJvdWdoIFwiZGVmYXVsdFZpZXcuZ2V0Q29tcHV0ZWRTdHlsZVwiXG5cdFx0dmFyIHZpZXcgPSBlbGVtLm93bmVyRG9jdW1lbnQuZGVmYXVsdFZpZXc7XG5cblx0XHRpZiAoICF2aWV3IHx8ICF2aWV3Lm9wZW5lciApIHtcblx0XHRcdHZpZXcgPSB3aW5kb3c7XG5cdFx0fVxuXG5cdFx0cmV0dXJuIHZpZXcuZ2V0Q29tcHV0ZWRTdHlsZSggZWxlbSApO1xuXHR9O1xufSApO1xuIl19
//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInd3dy9wYWNrcy9qcXVlcnkvc3JjL2Nzcy92YXIvZ2V0U3R5bGVzLmpzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBLE1BQU0sQ0FBQztJQUNILFlBQVksQ0FBQztJQUNiLE1BQU0sQ0FBQyxVQUFVLElBQUk7UUFDakIsdURBQXVEO1FBQ3ZELDBDQUEwQztRQUMxQywrRUFBK0U7UUFDL0UsSUFBSSxJQUFJLEdBQUcsSUFBSSxDQUFDLGFBQWEsQ0FBQyxXQUFXLENBQUM7UUFDMUMsRUFBRSxDQUFDLENBQUMsQ0FBQyxJQUFJLElBQUksQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLENBQUMsQ0FBQztZQUN4QixJQUFJLEdBQUcsTUFBTSxDQUFDO1FBQ2xCLENBQUM7UUFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLGdCQUFnQixDQUFDLElBQUksQ0FBQyxDQUFDO0lBQ3ZDLENBQUMsQ0FBQztBQUNOLENBQUMsQ0FBQyxDQUFDO0FBRUgsK3lDQUEreUMiLCJmaWxlIjoid3d3L3BhY2tzL2pxdWVyeS9zcmMvY3NzL3Zhci9nZXRTdHlsZXMuanMiLCJzb3VyY2VzQ29udGVudCI6WyJkZWZpbmUoZnVuY3Rpb24gKCkge1xuICAgIFwidXNlIHN0cmljdFwiO1xuICAgIHJldHVybiBmdW5jdGlvbiAoZWxlbSkge1xuICAgICAgICAvLyBTdXBwb3J0OiBJRSA8PTExIG9ubHksIEZpcmVmb3ggPD0zMCAoIzE1MDk4LCAjMTQxNTApXG4gICAgICAgIC8vIElFIHRocm93cyBvbiBlbGVtZW50cyBjcmVhdGVkIGluIHBvcHVwc1xuICAgICAgICAvLyBGRiBtZWFud2hpbGUgdGhyb3dzIG9uIGZyYW1lIGVsZW1lbnRzIHRocm91Z2ggXCJkZWZhdWx0Vmlldy5nZXRDb21wdXRlZFN0eWxlXCJcbiAgICAgICAgdmFyIHZpZXcgPSBlbGVtLm93bmVyRG9jdW1lbnQuZGVmYXVsdFZpZXc7XG4gICAgICAgIGlmICghdmlldyB8fCAhdmlldy5vcGVuZXIpIHtcbiAgICAgICAgICAgIHZpZXcgPSB3aW5kb3c7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIHZpZXcuZ2V0Q29tcHV0ZWRTdHlsZShlbGVtKTtcbiAgICB9O1xufSk7XG5cbi8vIyBzb3VyY2VNYXBwaW5nVVJMPWRhdGE6YXBwbGljYXRpb24vanNvbjtjaGFyc2V0PXV0Zjg7YmFzZTY0LGV5SjJaWEp6YVc5dUlqb3pMQ0p6YjNWeVkyVnpJanBiSW5kM2R5OXdZV05yY3k5cWNYVmxjbmt2YzNKakwyTnpjeTkyWVhJdloyVjBVM1I1YkdWekxtcHpJbDBzSW01aGJXVnpJanBiWFN3aWJXRndjR2x1WjNNaU9pSkJRVUZCTEUxQlFVMHNRMEZCUlR0SlFVTlFMRmxCUVZrc1EwRkJRenRKUVVWaUxFMUJRVTBzUTBGQlF5eFZRVUZWTEVsQlFVazdVVUZGY0VJc2RVUkJRWFZFTzFGQlEzWkVMREJEUVVFd1F6dFJRVU14UXl3clJVRkJLMFU3VVVGREwwVXNTVUZCU1N4SlFVRkpMRWRCUVVjc1NVRkJTU3hEUVVGRExHRkJRV0VzUTBGQlF5eFhRVUZYTEVOQlFVTTdVVUZGTVVNc1JVRkJSU3hEUVVGRExFTkJRVVVzUTBGQlF5eEpRVUZKTEVsQlFVa3NRMEZCUXl4SlFVRkpMRU5CUVVNc1RVRkJUeXhEUVVGRExFTkJRVU1zUTBGQlF6dFpRVU0zUWl4SlFVRkpMRWRCUVVjc1RVRkJUU3hEUVVGRE8xRkJRMllzUTBGQlF6dFJRVVZFTEUxQlFVMHNRMEZCUXl4SlFVRkpMRU5CUVVNc1owSkJRV2RDTEVOQlFVVXNTVUZCU1N4RFFVRkZMRU5CUVVNN1NVRkRkRU1zUTBGQlF5eERRVUZETzBGQlEwZ3NRMEZCUXl4RFFVRkZMRU5CUVVNaUxDSm1hV3hsSWpvaWQzZDNMM0JoWTJ0ekwycHhkV1Z5ZVM5emNtTXZZM056TDNaaGNpOW5aWFJUZEhsc1pYTXVhbk1pTENKemIzVnlZMlZ6UTI5dWRHVnVkQ0k2V3lKa1pXWnBibVVvSUdaMWJtTjBhVzl1S0NrZ2UxeHVYSFJjSW5WelpTQnpkSEpwWTNSY0lqdGNibHh1WEhSeVpYUjFjbTRnWm5WdVkzUnBiMjRvSUdWc1pXMGdLU0I3WEc1Y2JseDBYSFF2THlCVGRYQndiM0owT2lCSlJTQThQVEV4SUc5dWJIa3NJRVpwY21WbWIzZ2dQRDB6TUNBb0l6RTFNRGs0TENBak1UUXhOVEFwWEc1Y2RGeDBMeThnU1VVZ2RHaHliM2R6SUc5dUlHVnNaVzFsYm5SeklHTnlaV0YwWldRZ2FXNGdjRzl3ZFhCelhHNWNkRngwTHk4Z1JrWWdiV1ZoYm5kb2FXeGxJSFJvY205M2N5QnZiaUJtY21GdFpTQmxiR1Z0Wlc1MGN5QjBhSEp2ZFdkb0lGd2laR1ZtWVhWc2RGWnBaWGN1WjJWMFEyOXRjSFYwWldSVGRIbHNaVndpWEc1Y2RGeDBkbUZ5SUhacFpYY2dQU0JsYkdWdExtOTNibVZ5Ukc5amRXMWxiblF1WkdWbVlYVnNkRlpwWlhjN1hHNWNibHgwWEhScFppQW9JQ0YyYVdWM0lIeDhJQ0YyYVdWM0xtOXdaVzVsY2lBcElIdGNibHgwWEhSY2RIWnBaWGNnUFNCM2FXNWtiM2M3WEc1Y2RGeDBmVnh1WEc1Y2RGeDBjbVYwZFhKdUlIWnBaWGN1WjJWMFEyOXRjSFYwWldSVGRIbHNaU2dnWld4bGJTQXBPMXh1WEhSOU8xeHVmU0FwTzF4dUlsMTlcbiJdfQ==

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInd3dy9wYWNrcy9qcXVlcnkvc3JjL2Nzcy92YXIvZ2V0U3R5bGVzLmpzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQUFBLE1BQU0sQ0FBQztJQUNILFlBQVksQ0FBQztJQUNiLE1BQU0sQ0FBQyxVQUFVLElBQUk7UUFDakIsdURBQXVEO1FBQ3ZELDBDQUEwQztRQUMxQywrRUFBK0U7UUFDL0UsSUFBSSxJQUFJLEdBQUcsSUFBSSxDQUFDLGFBQWEsQ0FBQyxXQUFXLENBQUM7UUFDMUMsRUFBRSxDQUFDLENBQUMsQ0FBQyxJQUFJLElBQUksQ0FBQyxJQUFJLENBQUMsTUFBTSxDQUFDLENBQUMsQ0FBQztZQUN4QixJQUFJLEdBQUcsTUFBTSxDQUFDO1FBQ2xCLENBQUM7UUFDRCxNQUFNLENBQUMsSUFBSSxDQUFDLGdCQUFnQixDQUFDLElBQUksQ0FBQyxDQUFDO0lBQ3ZDLENBQUMsQ0FBQztBQUNOLENBQUMsQ0FBQyxDQUFDO0FBQ0gsK3lDQUEreUM7QUFFL3lDLHVsR0FBdWxHIiwiZmlsZSI6Ind3dy9wYWNrcy9qcXVlcnkvc3JjL2Nzcy92YXIvZ2V0U3R5bGVzLmpzIiwic291cmNlc0NvbnRlbnQiOlsiZGVmaW5lKGZ1bmN0aW9uICgpIHtcbiAgICBcInVzZSBzdHJpY3RcIjtcbiAgICByZXR1cm4gZnVuY3Rpb24gKGVsZW0pIHtcbiAgICAgICAgLy8gU3VwcG9ydDogSUUgPD0xMSBvbmx5LCBGaXJlZm94IDw9MzAgKCMxNTA5OCwgIzE0MTUwKVxuICAgICAgICAvLyBJRSB0aHJvd3Mgb24gZWxlbWVudHMgY3JlYXRlZCBpbiBwb3B1cHNcbiAgICAgICAgLy8gRkYgbWVhbndoaWxlIHRocm93cyBvbiBmcmFtZSBlbGVtZW50cyB0aHJvdWdoIFwiZGVmYXVsdFZpZXcuZ2V0Q29tcHV0ZWRTdHlsZVwiXG4gICAgICAgIHZhciB2aWV3ID0gZWxlbS5vd25lckRvY3VtZW50LmRlZmF1bHRWaWV3O1xuICAgICAgICBpZiAoIXZpZXcgfHwgIXZpZXcub3BlbmVyKSB7XG4gICAgICAgICAgICB2aWV3ID0gd2luZG93O1xuICAgICAgICB9XG4gICAgICAgIHJldHVybiB2aWV3LmdldENvbXB1dGVkU3R5bGUoZWxlbSk7XG4gICAgfTtcbn0pO1xuLy8jIHNvdXJjZU1hcHBpbmdVUkw9ZGF0YTphcHBsaWNhdGlvbi9qc29uO2NoYXJzZXQ9dXRmODtiYXNlNjQsZXlKMlpYSnphVzl1SWpvekxDSnpiM1Z5WTJWeklqcGJJbmQzZHk5d1lXTnJjeTlxY1hWbGNua3ZjM0pqTDJOemN5OTJZWEl2WjJWMFUzUjViR1Z6TG1weklsMHNJbTVoYldWeklqcGJYU3dpYldGd2NHbHVaM01pT2lKQlFVRkJMRTFCUVUwc1EwRkJSVHRKUVVOUUxGbEJRVmtzUTBGQlF6dEpRVVZpTEUxQlFVMHNRMEZCUXl4VlFVRlZMRWxCUVVrN1VVRkZjRUlzZFVSQlFYVkVPMUZCUTNaRUxEQkRRVUV3UXp0UlFVTXhReXdyUlVGQkswVTdVVUZETDBVc1NVRkJTU3hKUVVGSkxFZEJRVWNzU1VGQlNTeERRVUZETEdGQlFXRXNRMEZCUXl4WFFVRlhMRU5CUVVNN1VVRkZNVU1zUlVGQlJTeERRVUZETEVOQlFVVXNRMEZCUXl4SlFVRkpMRWxCUVVrc1EwRkJReXhKUVVGSkxFTkJRVU1zVFVGQlR5eERRVUZETEVOQlFVTXNRMEZCUXp0WlFVTTNRaXhKUVVGSkxFZEJRVWNzVFVGQlRTeERRVUZETzFGQlEyWXNRMEZCUXp0UlFVVkVMRTFCUVUwc1EwRkJReXhKUVVGSkxFTkJRVU1zWjBKQlFXZENMRU5CUVVVc1NVRkJTU3hEUVVGRkxFTkJRVU03U1VGRGRFTXNRMEZCUXl4RFFVRkRPMEZCUTBnc1EwRkJReXhEUVVGRkxFTkJRVU1pTENKbWFXeGxJam9pZDNkM0wzQmhZMnR6TDJweGRXVnllUzl6Y21NdlkzTnpMM1poY2k5blpYUlRkSGxzWlhNdWFuTWlMQ0p6YjNWeVkyVnpRMjl1ZEdWdWRDSTZXeUprWldacGJtVW9JR1oxYm1OMGFXOXVLQ2tnZTF4dVhIUmNJblZ6WlNCemRISnBZM1JjSWp0Y2JseHVYSFJ5WlhSMWNtNGdablZ1WTNScGIyNG9JR1ZzWlcwZ0tTQjdYRzVjYmx4MFhIUXZMeUJUZFhCd2IzSjBPaUJKUlNBOFBURXhJRzl1Ykhrc0lFWnBjbVZtYjNnZ1BEMHpNQ0FvSXpFMU1EazRMQ0FqTVRReE5UQXBYRzVjZEZ4MEx5OGdTVVVnZEdoeWIzZHpJRzl1SUdWc1pXMWxiblJ6SUdOeVpXRjBaV1FnYVc0Z2NHOXdkWEJ6WEc1Y2RGeDBMeThnUmtZZ2JXVmhibmRvYVd4bElIUm9jbTkzY3lCdmJpQm1jbUZ0WlNCbGJHVnRaVzUwY3lCMGFISnZkV2RvSUZ3aVpHVm1ZWFZzZEZacFpYY3VaMlYwUTI5dGNIVjBaV1JUZEhsc1pWd2lYRzVjZEZ4MGRtRnlJSFpwWlhjZ1BTQmxiR1Z0TG05M2JtVnlSRzlqZFcxbGJuUXVaR1ZtWVhWc2RGWnBaWGM3WEc1Y2JseDBYSFJwWmlBb0lDRjJhV1YzSUh4OElDRjJhV1YzTG05d1pXNWxjaUFwSUh0Y2JseDBYSFJjZEhacFpYY2dQU0IzYVc1a2IzYzdYRzVjZEZ4MGZWeHVYRzVjZEZ4MGNtVjBkWEp1SUhacFpYY3VaMlYwUTI5dGNIVjBaV1JUZEhsc1pTZ2daV3hsYlNBcE8xeHVYSFI5TzF4dWZTQXBPMXh1SWwxOVxuXG4vLyMgc291cmNlTWFwcGluZ1VSTD1kYXRhOmFwcGxpY2F0aW9uL2pzb247Y2hhcnNldD11dGY4O2Jhc2U2NCxleUoyWlhKemFXOXVJam96TENKemIzVnlZMlZ6SWpwYkluZDNkeTl3WVdOcmN5OXFjWFZsY25rdmMzSmpMMk56Y3k5MllYSXZaMlYwVTNSNWJHVnpMbXB6SWwwc0ltNWhiV1Z6SWpwYlhTd2liV0Z3Y0dsdVozTWlPaUpCUVVGQkxFMUJRVTBzUTBGQlF6dEpRVU5JTEZsQlFWa3NRMEZCUXp0SlFVTmlMRTFCUVUwc1EwRkJReXhWUVVGVkxFbEJRVWs3VVVGRGFrSXNkVVJCUVhWRU8xRkJRM1pFTERCRFFVRXdRenRSUVVNeFF5d3JSVUZCSzBVN1VVRkRMMFVzU1VGQlNTeEpRVUZKTEVkQlFVY3NTVUZCU1N4RFFVRkRMR0ZCUVdFc1EwRkJReXhYUVVGWExFTkJRVU03VVVGRE1VTXNSVUZCUlN4RFFVRkRMRU5CUVVNc1EwRkJReXhKUVVGSkxFbEJRVWtzUTBGQlF5eEpRVUZKTEVOQlFVTXNUVUZCVFN4RFFVRkRMRU5CUVVNc1EwRkJRenRaUVVONFFpeEpRVUZKTEVkQlFVY3NUVUZCVFN4RFFVRkRPMUZCUTJ4Q0xFTkJRVU03VVVGRFJDeE5RVUZOTEVOQlFVTXNTVUZCU1N4RFFVRkRMR2RDUVVGblFpeERRVUZETEVsQlFVa3NRMEZCUXl4RFFVRkRPMGxCUTNaRExFTkJRVU1zUTBGQlF6dEJRVU5PTEVOQlFVTXNRMEZCUXl4RFFVRkRPMEZCUlVnc0szbERRVUVyZVVNaUxDSm1hV3hsSWpvaWQzZDNMM0JoWTJ0ekwycHhkV1Z5ZVM5emNtTXZZM056TDNaaGNpOW5aWFJUZEhsc1pYTXVhbk1pTENKemIzVnlZMlZ6UTI5dWRHVnVkQ0k2V3lKa1pXWnBibVVvWm5WdVkzUnBiMjRnS0NrZ2UxeHVJQ0FnSUZ3aWRYTmxJSE4wY21samRGd2lPMXh1SUNBZ0lISmxkSFZ5YmlCbWRXNWpkR2x2YmlBb1pXeGxiU2tnZTF4dUlDQWdJQ0FnSUNBdkx5QlRkWEJ3YjNKME9pQkpSU0E4UFRFeElHOXViSGtzSUVacGNtVm1iM2dnUEQwek1DQW9JekUxTURrNExDQWpNVFF4TlRBcFhHNGdJQ0FnSUNBZ0lDOHZJRWxGSUhSb2NtOTNjeUJ2YmlCbGJHVnRaVzUwY3lCamNtVmhkR1ZrSUdsdUlIQnZjSFZ3YzF4dUlDQWdJQ0FnSUNBdkx5QkdSaUJ0WldGdWQyaHBiR1VnZEdoeWIzZHpJRzl1SUdaeVlXMWxJR1ZzWlcxbGJuUnpJSFJvY205MVoyZ2dYQ0prWldaaGRXeDBWbWxsZHk1blpYUkRiMjF3ZFhSbFpGTjBlV3hsWENKY2JpQWdJQ0FnSUNBZ2RtRnlJSFpwWlhjZ1BTQmxiR1Z0TG05M2JtVnlSRzlqZFcxbGJuUXVaR1ZtWVhWc2RGWnBaWGM3WEc0Z0lDQWdJQ0FnSUdsbUlDZ2hkbWxsZHlCOGZDQWhkbWxsZHk1dmNHVnVaWElwSUh0Y2JpQWdJQ0FnSUNBZ0lDQWdJSFpwWlhjZ1BTQjNhVzVrYjNjN1hHNGdJQ0FnSUNBZ0lIMWNiaUFnSUNBZ0lDQWdjbVYwZFhKdUlIWnBaWGN1WjJWMFEyOXRjSFYwWldSVGRIbHNaU2hsYkdWdEtUdGNiaUFnSUNCOU8xeHVmU2s3WEc1Y2JpOHZJeUJ6YjNWeVkyVk5ZWEJ3YVc1blZWSk1QV1JoZEdFNllYQndiR2xqWVhScGIyNHZhbk52Ymp0amFHRnljMlYwUFhWMFpqZzdZbUZ6WlRZMExHVjVTakphV0VwNllWYzVkVWxxYjNwTVEwcDZZak5XZVZreVZucEphbkJpU1c1a00yUjVPWGRaVjA1eVkzazVjV05ZVm14amJtdDJZek5LYWt3eVRucGplVGt5V1ZoSmRsb3lWakJWTTFJMVlrZFdla3h0Y0hwSmJEQnpTVzAxYUdKWFZucEphbkJpV0ZOM2FXSlhSbmRqUjJ4MVdqTk5hVTlwU2tKUlZVWkNURVV4UWxGVk1ITlJNRVpDVWxSMFNsRlZUbEZNUm14Q1VWWnJjMUV3UmtKUmVuUktVVlZXYVV4Rk1VSlJWVEJ6VVRCR1FsRjVlRlpSVlVaV1RFVnNRbEZWYXpkVlZVWkdZMFZKYzJSVlVrSlJXRlpGVHpGR1FsRXpXa1ZNUkVKRVVWVkZkMUY2ZEZKUlZVMTRVWGwzY2xKVlJrSkxNRlUzVlZWR1JFd3dWWE5UVlVaQ1UxTjRTbEZWUmtwTVJXUkNVVlZqYzFOVlJrSlRVM2hFVVZWR1JFeEhSa0pSVjBWelVUQkdRbEY1ZUZoUlZVWllURVZPUWxGVlRUZFZWVVpHVFZWTmMxSlZSa0pTVTNoRVVWVkdSRXhGVGtKUlZWVnpVVEJHUWxGNWVFcFJWVVpLVEVWc1FsRlZhM05STUVaQ1VYbDRTbEZWUmtwTVJVNUNVVlZOYzFSVlJrSlVlWGhFVVZWR1JFeEZUa0pSVlUxelVUQkdRbEY2ZEZwUlZVMHpVV2w0U2xGVlJrcE1SV1JDVVZWamMxUlZSa0pVVTNoRVVWVkdSRTh4UmtKUk1sbHpVVEJHUWxGNmRGSlJWVlpGVEVVeFFsRlZNSE5STUVaQ1VYbDRTbEZWUmtwTVJVNUNVVlZOYzFvd1NrSlJWMlJEVEVWT1FsRlZWWE5UVlVaQ1UxTjRSRkZWUmtaTVJVNUNVVlZOTjFOVlJrUmtSVTF6VVRCR1FsRjVlRVJSVlVaRVR6QkdRbEV3WjNOUk1FWkNVWGw0UkZGVlJrWk1SVTVDVVZWTmFVeERTbTFoVjNoc1NXcHZhV1F6WkROTU0wSm9XVEowZWt3eWNIaGtWMVo1WlZNNWVtTnRUWFpaTTA1NlRETmFhR05wT1c1YVdGSlVaRWhzYzFwWVRYVmhiazFwVEVOS2VtSXpWbmxaTWxaNlVUSTVkV1JIVm5Wa1EwazJWM2xLYTFwWFduQmliVlZ2U1VkYU1XSnRUakJoVnpsMVMwTnJaMlV4ZUhWWVNGSmpTVzVXZWxwVFFucGtTRXB3V1ROU1kwbHFkR05pYkhoMVdFaFNlVnBZVWpGamJUUm5XbTVXZFZrelVuQmlNalJ2U1VkV2MxcFhNR2RMVTBJM1dFYzFZMkpzZURCWVNGRjJUSGxDVkdSWVFuZGlNMG93VDJsQ1NsSlRRVGhRVkVWNFNVYzVkV0pJYTNOSlJWcHdZMjFXYldJeloyZFFSREI2VFVOQmIwbDZSVEZOUkdzMFRFTkJhazFVVVhoT1ZFRndXRWMxWTJSR2VEQk1lVGhuVTFWVloyUkhhSGxpTTJSNlNVYzVkVWxIVm5OYVZ6RnNZbTVTZWtsSFRubGFWMFl3V2xkUloyRlhOR2RqUnpsM1pGaENlbGhITldOa1JuZ3dUSGs0WjFKcldXZGlWMVpvWW01a2IyRlhlR3hKU0ZKdlkyMDVNMk41UW5aaWFVSnRZMjFHZEZwVFFteGlSMVowV2xjMU1HTjVRakJoU0VwMlpGZGtiMGxHZDJsYVIxWnRXVmhXYzJSR1duQmFXR04xV2pKV01GRXlPWFJqU0ZZd1dsZFNWR1JJYkhOYVZuZHBXRWMxWTJSR2VEQmtiVVo1U1VoYWNGcFlZMmRRVTBKc1lrZFdkRXh0T1ROaWJWWjVVa2M1YW1SWE1XeGlibEYxV2tkV2JWbFlWbk5rUmxwd1dsaGpOMWhITldOaWJIZ3dXRWhTY0ZwcFFXOUpRMFl5WVZkV00wbEllRGhKUTBZeVlWZFdNMHh0T1hkYVZ6VnNZMmxCY0VsSWRHTmliSGd3V0VoU1kyUklXbkJhV0dOblVGTkNNMkZYTld0aU0yTTNXRWMxWTJSR2VEQm1WbmgxV0VjMVkyUkdlREJqYlZZd1pGaEtkVWxJV25CYVdHTjFXakpXTUZFeU9YUmpTRll3V2xkU1ZHUkliSE5hVTJkbldsZDRiR0pUUVhCUE1YaDFXRWhTT1U4eGVIVm1VMEZ3VHpGNGRVbHNNVGxjYmlKZGZRPT1cbiJdfQ==