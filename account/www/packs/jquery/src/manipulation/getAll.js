define([
    "../core",
    "../core/nodeName"
], function (jQuery, nodeName) {
    "use strict";
    function getAll(context, tag) {
        // Support: IE <=9 - 11 only
        // Use typeof to avoid zero-argument method invocation on host objects (#15151)
        var ret;
        if (typeof context.getElementsByTagName !== "undefined") {
            ret = context.getElementsByTagName(tag || "*");
        }
        else if (typeof context.querySelectorAll !== "undefined") {
            ret = context.querySelectorAll(tag || "*");
        }
        else {
            ret = [];
        }
        if (tag === undefined || tag && nodeName(context, tag)) {
            return jQuery.merge([context], ret);
        }
        return ret;
    }
    return getAll;
});
//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInd3dy9wYWNrcy9qcXVlcnkvc3JjL21hbmlwdWxhdGlvbi9nZXRBbGwuanMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUEsTUFBTSxDQUFFO0lBQ1AsU0FBUztJQUNULGtCQUFrQjtDQUNsQixFQUFFLFVBQVUsTUFBTSxFQUFFLFFBQVE7SUFFN0IsWUFBWSxDQUFDO0lBRWIsZ0JBQWlCLE9BQU8sRUFBRSxHQUFHO1FBRTVCLDRCQUE0QjtRQUM1QiwrRUFBK0U7UUFDL0UsSUFBSSxHQUFHLENBQUM7UUFFUixFQUFFLENBQUMsQ0FBRSxPQUFPLE9BQU8sQ0FBQyxvQkFBb0IsS0FBSyxXQUFZLENBQUMsQ0FBQyxDQUFDO1lBQzNELEdBQUcsR0FBRyxPQUFPLENBQUMsb0JBQW9CLENBQUUsR0FBRyxJQUFJLEdBQUcsQ0FBRSxDQUFDO1FBRWxELENBQUM7UUFBQyxJQUFJLENBQUMsRUFBRSxDQUFDLENBQUUsT0FBTyxPQUFPLENBQUMsZ0JBQWdCLEtBQUssV0FBWSxDQUFDLENBQUMsQ0FBQztZQUM5RCxHQUFHLEdBQUcsT0FBTyxDQUFDLGdCQUFnQixDQUFFLEdBQUcsSUFBSSxHQUFHLENBQUUsQ0FBQztRQUU5QyxDQUFDO1FBQUMsSUFBSSxDQUFDLENBQUM7WUFDUCxHQUFHLEdBQUcsRUFBRSxDQUFDO1FBQ1YsQ0FBQztRQUVELEVBQUUsQ0FBQyxDQUFFLEdBQUcsS0FBSyxTQUFTLElBQUksR0FBRyxJQUFJLFFBQVEsQ0FBRSxPQUFPLEVBQUUsR0FBRyxDQUFHLENBQUMsQ0FBQyxDQUFDO1lBQzVELE1BQU0sQ0FBQyxNQUFNLENBQUMsS0FBSyxDQUFFLENBQUUsT0FBTyxDQUFFLEVBQUUsR0FBRyxDQUFFLENBQUM7UUFDekMsQ0FBQztRQUVELE1BQU0sQ0FBQyxHQUFHLENBQUM7SUFDWixDQUFDO0lBRUQsTUFBTSxDQUFDLE1BQU0sQ0FBQztBQUNkLENBQUMsQ0FBRSxDQUFDIiwiZmlsZSI6Ind3dy9wYWNrcy9qcXVlcnkvc3JjL21hbmlwdWxhdGlvbi9nZXRBbGwuanMiLCJzb3VyY2VzQ29udGVudCI6WyJkZWZpbmUoIFtcblx0XCIuLi9jb3JlXCIsXG5cdFwiLi4vY29yZS9ub2RlTmFtZVwiXG5dLCBmdW5jdGlvbiggalF1ZXJ5LCBub2RlTmFtZSApIHtcblxuXCJ1c2Ugc3RyaWN0XCI7XG5cbmZ1bmN0aW9uIGdldEFsbCggY29udGV4dCwgdGFnICkge1xuXG5cdC8vIFN1cHBvcnQ6IElFIDw9OSAtIDExIG9ubHlcblx0Ly8gVXNlIHR5cGVvZiB0byBhdm9pZCB6ZXJvLWFyZ3VtZW50IG1ldGhvZCBpbnZvY2F0aW9uIG9uIGhvc3Qgb2JqZWN0cyAoIzE1MTUxKVxuXHR2YXIgcmV0O1xuXG5cdGlmICggdHlwZW9mIGNvbnRleHQuZ2V0RWxlbWVudHNCeVRhZ05hbWUgIT09IFwidW5kZWZpbmVkXCIgKSB7XG5cdFx0cmV0ID0gY29udGV4dC5nZXRFbGVtZW50c0J5VGFnTmFtZSggdGFnIHx8IFwiKlwiICk7XG5cblx0fSBlbHNlIGlmICggdHlwZW9mIGNvbnRleHQucXVlcnlTZWxlY3RvckFsbCAhPT0gXCJ1bmRlZmluZWRcIiApIHtcblx0XHRyZXQgPSBjb250ZXh0LnF1ZXJ5U2VsZWN0b3JBbGwoIHRhZyB8fCBcIipcIiApO1xuXG5cdH0gZWxzZSB7XG5cdFx0cmV0ID0gW107XG5cdH1cblxuXHRpZiAoIHRhZyA9PT0gdW5kZWZpbmVkIHx8IHRhZyAmJiBub2RlTmFtZSggY29udGV4dCwgdGFnICkgKSB7XG5cdFx0cmV0dXJuIGpRdWVyeS5tZXJnZSggWyBjb250ZXh0IF0sIHJldCApO1xuXHR9XG5cblx0cmV0dXJuIHJldDtcbn1cblxucmV0dXJuIGdldEFsbDtcbn0gKTtcbiJdfQ==
//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInd3dy9wYWNrcy9qcXVlcnkvc3JjL21hbmlwdWxhdGlvbi9nZXRBbGwuanMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUEsTUFBTSxDQUFDO0lBQ0gsU0FBUztJQUNULGtCQUFrQjtDQUNyQixFQUFFLFVBQVUsTUFBTSxFQUFFLFFBQVE7SUFDekIsWUFBWSxDQUFDO0lBQ2IsZ0JBQWdCLE9BQU8sRUFBRSxHQUFHO1FBQ3hCLDRCQUE0QjtRQUM1QiwrRUFBK0U7UUFDL0UsSUFBSSxHQUFHLENBQUM7UUFDUixFQUFFLENBQUMsQ0FBQyxPQUFPLE9BQU8sQ0FBQyxvQkFBb0IsS0FBSyxXQUFXLENBQUMsQ0FBQyxDQUFDO1lBQ3RELEdBQUcsR0FBRyxPQUFPLENBQUMsb0JBQW9CLENBQUMsR0FBRyxJQUFJLEdBQUcsQ0FBQyxDQUFDO1FBQ25ELENBQUM7UUFDRCxJQUFJLENBQUMsRUFBRSxDQUFDLENBQUMsT0FBTyxPQUFPLENBQUMsZ0JBQWdCLEtBQUssV0FBVyxDQUFDLENBQUMsQ0FBQztZQUN2RCxHQUFHLEdBQUcsT0FBTyxDQUFDLGdCQUFnQixDQUFDLEdBQUcsSUFBSSxHQUFHLENBQUMsQ0FBQztRQUMvQyxDQUFDO1FBQ0QsSUFBSSxDQUFDLENBQUM7WUFDRixHQUFHLEdBQUcsRUFBRSxDQUFDO1FBQ2IsQ0FBQztRQUNELEVBQUUsQ0FBQyxDQUFDLEdBQUcsS0FBSyxTQUFTLElBQUksR0FBRyxJQUFJLFFBQVEsQ0FBQyxPQUFPLEVBQUUsR0FBRyxDQUFDLENBQUMsQ0FBQyxDQUFDO1lBQ3JELE1BQU0sQ0FBQyxNQUFNLENBQUMsS0FBSyxDQUFDLENBQUMsT0FBTyxDQUFDLEVBQUUsR0FBRyxDQUFDLENBQUM7UUFDeEMsQ0FBQztRQUNELE1BQU0sQ0FBQyxHQUFHLENBQUM7SUFDZixDQUFDO0lBQ0QsTUFBTSxDQUFDLE1BQU0sQ0FBQztBQUNsQixDQUFDLENBQUMsQ0FBQztBQUVILHV0RUFBdXRFIiwiZmlsZSI6Ind3dy9wYWNrcy9qcXVlcnkvc3JjL21hbmlwdWxhdGlvbi9nZXRBbGwuanMiLCJzb3VyY2VzQ29udGVudCI6WyJkZWZpbmUoW1xuICAgIFwiLi4vY29yZVwiLFxuICAgIFwiLi4vY29yZS9ub2RlTmFtZVwiXG5dLCBmdW5jdGlvbiAoalF1ZXJ5LCBub2RlTmFtZSkge1xuICAgIFwidXNlIHN0cmljdFwiO1xuICAgIGZ1bmN0aW9uIGdldEFsbChjb250ZXh0LCB0YWcpIHtcbiAgICAgICAgLy8gU3VwcG9ydDogSUUgPD05IC0gMTEgb25seVxuICAgICAgICAvLyBVc2UgdHlwZW9mIHRvIGF2b2lkIHplcm8tYXJndW1lbnQgbWV0aG9kIGludm9jYXRpb24gb24gaG9zdCBvYmplY3RzICgjMTUxNTEpXG4gICAgICAgIHZhciByZXQ7XG4gICAgICAgIGlmICh0eXBlb2YgY29udGV4dC5nZXRFbGVtZW50c0J5VGFnTmFtZSAhPT0gXCJ1bmRlZmluZWRcIikge1xuICAgICAgICAgICAgcmV0ID0gY29udGV4dC5nZXRFbGVtZW50c0J5VGFnTmFtZSh0YWcgfHwgXCIqXCIpO1xuICAgICAgICB9XG4gICAgICAgIGVsc2UgaWYgKHR5cGVvZiBjb250ZXh0LnF1ZXJ5U2VsZWN0b3JBbGwgIT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgICAgICAgIHJldCA9IGNvbnRleHQucXVlcnlTZWxlY3RvckFsbCh0YWcgfHwgXCIqXCIpO1xuICAgICAgICB9XG4gICAgICAgIGVsc2Uge1xuICAgICAgICAgICAgcmV0ID0gW107XG4gICAgICAgIH1cbiAgICAgICAgaWYgKHRhZyA9PT0gdW5kZWZpbmVkIHx8IHRhZyAmJiBub2RlTmFtZShjb250ZXh0LCB0YWcpKSB7XG4gICAgICAgICAgICByZXR1cm4galF1ZXJ5Lm1lcmdlKFtjb250ZXh0XSwgcmV0KTtcbiAgICAgICAgfVxuICAgICAgICByZXR1cm4gcmV0O1xuICAgIH1cbiAgICByZXR1cm4gZ2V0QWxsO1xufSk7XG5cbi8vIyBzb3VyY2VNYXBwaW5nVVJMPWRhdGE6YXBwbGljYXRpb24vanNvbjtjaGFyc2V0PXV0Zjg7YmFzZTY0LGV5SjJaWEp6YVc5dUlqb3pMQ0p6YjNWeVkyVnpJanBiSW5kM2R5OXdZV05yY3k5cWNYVmxjbmt2YzNKakwyMWhibWx3ZFd4aGRHbHZiaTluWlhSQmJHd3Vhbk1pWFN3aWJtRnRaWE1pT2x0ZExDSnRZWEJ3YVc1bmN5STZJa0ZCUVVFc1RVRkJUU3hEUVVGRk8wbEJRMUFzVTBGQlV6dEpRVU5VTEd0Q1FVRnJRanREUVVOc1FpeEZRVUZGTEZWQlFWVXNUVUZCVFN4RlFVRkZMRkZCUVZFN1NVRkZOMElzV1VGQldTeERRVUZETzBsQlJXSXNaMEpCUVdsQ0xFOUJRVThzUlVGQlJTeEhRVUZITzFGQlJUVkNMRFJDUVVFMFFqdFJRVU0xUWl3clJVRkJLMFU3VVVGREwwVXNTVUZCU1N4SFFVRkhMRU5CUVVNN1VVRkZVaXhGUVVGRkxFTkJRVU1zUTBGQlJTeFBRVUZQTEU5QlFVOHNRMEZCUXl4dlFrRkJiMElzUzBGQlN5eFhRVUZaTEVOQlFVTXNRMEZCUXl4RFFVRkRPMWxCUXpORUxFZEJRVWNzUjBGQlJ5eFBRVUZQTEVOQlFVTXNiMEpCUVc5Q0xFTkJRVVVzUjBGQlJ5eEpRVUZKTEVkQlFVY3NRMEZCUlN4RFFVRkRPMUZCUld4RUxFTkJRVU03VVVGQlF5eEpRVUZKTEVOQlFVTXNSVUZCUlN4RFFVRkRMRU5CUVVVc1QwRkJUeXhQUVVGUExFTkJRVU1zWjBKQlFXZENMRXRCUVVzc1YwRkJXU3hEUVVGRExFTkJRVU1zUTBGQlF6dFpRVU01UkN4SFFVRkhMRWRCUVVjc1QwRkJUeXhEUVVGRExHZENRVUZuUWl4RFFVRkZMRWRCUVVjc1NVRkJTU3hIUVVGSExFTkJRVVVzUTBGQlF6dFJRVVU1UXl4RFFVRkRPMUZCUVVNc1NVRkJTU3hEUVVGRExFTkJRVU03V1VGRFVDeEhRVUZITEVkQlFVY3NSVUZCUlN4RFFVRkRPMUZCUTFZc1EwRkJRenRSUVVWRUxFVkJRVVVzUTBGQlF5eERRVUZGTEVkQlFVY3NTMEZCU3l4VFFVRlRMRWxCUVVrc1IwRkJSeXhKUVVGSkxGRkJRVkVzUTBGQlJTeFBRVUZQTEVWQlFVVXNSMEZCUnl4RFFVRkhMRU5CUVVNc1EwRkJReXhEUVVGRE8xbEJRelZFTEUxQlFVMHNRMEZCUXl4TlFVRk5MRU5CUVVNc1MwRkJTeXhEUVVGRkxFTkJRVVVzVDBGQlR5eERRVUZGTEVWQlFVVXNSMEZCUnl4RFFVRkZMRU5CUVVNN1VVRkRla01zUTBGQlF6dFJRVVZFTEUxQlFVMHNRMEZCUXl4SFFVRkhMRU5CUVVNN1NVRkRXaXhEUVVGRE8wbEJSVVFzVFVGQlRTeERRVUZETEUxQlFVMHNRMEZCUXp0QlFVTmtMRU5CUVVNc1EwRkJSU3hEUVVGRElpd2labWxzWlNJNkluZDNkeTl3WVdOcmN5OXFjWFZsY25rdmMzSmpMMjFoYm1sd2RXeGhkR2x2Ymk5blpYUkJiR3d1YW5NaUxDSnpiM1Z5WTJWelEyOXVkR1Z1ZENJNld5SmtaV1pwYm1Vb0lGdGNibHgwWENJdUxpOWpiM0psWENJc1hHNWNkRndpTGk0dlkyOXlaUzl1YjJSbFRtRnRaVndpWEc1ZExDQm1kVzVqZEdsdmJpZ2dhbEYxWlhKNUxDQnViMlJsVG1GdFpTQXBJSHRjYmx4dVhDSjFjMlVnYzNSeWFXTjBYQ0k3WEc1Y2JtWjFibU4wYVc5dUlHZGxkRUZzYkNnZ1kyOXVkR1Y0ZEN3Z2RHRm5JQ2tnZTF4dVhHNWNkQzh2SUZOMWNIQnZjblE2SUVsRklEdzlPU0F0SURFeElHOXViSGxjYmx4MEx5OGdWWE5sSUhSNWNHVnZaaUIwYnlCaGRtOXBaQ0I2WlhKdkxXRnlaM1Z0Wlc1MElHMWxkR2h2WkNCcGJuWnZZMkYwYVc5dUlHOXVJR2h2YzNRZ2IySnFaV04wY3lBb0l6RTFNVFV4S1Z4dVhIUjJZWElnY21WME8xeHVYRzVjZEdsbUlDZ2dkSGx3Wlc5bUlHTnZiblJsZUhRdVoyVjBSV3hsYldWdWRITkNlVlJoWjA1aGJXVWdJVDA5SUZ3aWRXNWtaV1pwYm1Wa1hDSWdLU0I3WEc1Y2RGeDBjbVYwSUQwZ1kyOXVkR1Y0ZEM1blpYUkZiR1Z0Wlc1MGMwSjVWR0ZuVG1GdFpTZ2dkR0ZuSUh4OElGd2lLbHdpSUNrN1hHNWNibHgwZlNCbGJITmxJR2xtSUNnZ2RIbHdaVzltSUdOdmJuUmxlSFF1Y1hWbGNubFRaV3hsWTNSdmNrRnNiQ0FoUFQwZ1hDSjFibVJsWm1sdVpXUmNJaUFwSUh0Y2JseDBYSFJ5WlhRZ1BTQmpiMjUwWlhoMExuRjFaWEo1VTJWc1pXTjBiM0pCYkd3b0lIUmhaeUI4ZkNCY0lpcGNJaUFwTzF4dVhHNWNkSDBnWld4elpTQjdYRzVjZEZ4MGNtVjBJRDBnVzEwN1hHNWNkSDFjYmx4dVhIUnBaaUFvSUhSaFp5QTlQVDBnZFc1a1pXWnBibVZrSUh4OElIUmhaeUFtSmlCdWIyUmxUbUZ0WlNnZ1kyOXVkR1Y0ZEN3Z2RHRm5JQ2tnS1NCN1hHNWNkRngwY21WMGRYSnVJR3BSZFdWeWVTNXRaWEpuWlNnZ1d5QmpiMjUwWlhoMElGMHNJSEpsZENBcE8xeHVYSFI5WEc1Y2JseDBjbVYwZFhKdUlISmxkRHRjYm4xY2JseHVjbVYwZFhKdUlHZGxkRUZzYkR0Y2JuMGdLVHRjYmlKZGZRPT1cbiJdfQ==

//# sourceMappingURL=data:application/json;charset=utf8;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbInd3dy9wYWNrcy9qcXVlcnkvc3JjL21hbmlwdWxhdGlvbi9nZXRBbGwuanMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBQUEsTUFBTSxDQUFDO0lBQ0gsU0FBUztJQUNULGtCQUFrQjtDQUNyQixFQUFFLFVBQVUsTUFBTSxFQUFFLFFBQVE7SUFDekIsWUFBWSxDQUFDO0lBQ2IsZ0JBQWdCLE9BQU8sRUFBRSxHQUFHO1FBQ3hCLDRCQUE0QjtRQUM1QiwrRUFBK0U7UUFDL0UsSUFBSSxHQUFHLENBQUM7UUFDUixFQUFFLENBQUMsQ0FBQyxPQUFPLE9BQU8sQ0FBQyxvQkFBb0IsS0FBSyxXQUFXLENBQUMsQ0FBQyxDQUFDO1lBQ3RELEdBQUcsR0FBRyxPQUFPLENBQUMsb0JBQW9CLENBQUMsR0FBRyxJQUFJLEdBQUcsQ0FBQyxDQUFDO1FBQ25ELENBQUM7UUFDRCxJQUFJLENBQUMsRUFBRSxDQUFDLENBQUMsT0FBTyxPQUFPLENBQUMsZ0JBQWdCLEtBQUssV0FBVyxDQUFDLENBQUMsQ0FBQztZQUN2RCxHQUFHLEdBQUcsT0FBTyxDQUFDLGdCQUFnQixDQUFDLEdBQUcsSUFBSSxHQUFHLENBQUMsQ0FBQztRQUMvQyxDQUFDO1FBQ0QsSUFBSSxDQUFDLENBQUM7WUFDRixHQUFHLEdBQUcsRUFBRSxDQUFDO1FBQ2IsQ0FBQztRQUNELEVBQUUsQ0FBQyxDQUFDLEdBQUcsS0FBSyxTQUFTLElBQUksR0FBRyxJQUFJLFFBQVEsQ0FBQyxPQUFPLEVBQUUsR0FBRyxDQUFDLENBQUMsQ0FBQyxDQUFDO1lBQ3JELE1BQU0sQ0FBQyxNQUFNLENBQUMsS0FBSyxDQUFDLENBQUMsT0FBTyxDQUFDLEVBQUUsR0FBRyxDQUFDLENBQUM7UUFDeEMsQ0FBQztRQUNELE1BQU0sQ0FBQyxHQUFHLENBQUM7SUFDZixDQUFDO0lBQ0QsTUFBTSxDQUFDLE1BQU0sQ0FBQztBQUNsQixDQUFDLENBQUMsQ0FBQztBQUNILHV0RUFBdXRFO0FBRXZ0RSxtektBQW16SyIsImZpbGUiOiJ3d3cvcGFja3MvanF1ZXJ5L3NyYy9tYW5pcHVsYXRpb24vZ2V0QWxsLmpzIiwic291cmNlc0NvbnRlbnQiOlsiZGVmaW5lKFtcbiAgICBcIi4uL2NvcmVcIixcbiAgICBcIi4uL2NvcmUvbm9kZU5hbWVcIlxuXSwgZnVuY3Rpb24gKGpRdWVyeSwgbm9kZU5hbWUpIHtcbiAgICBcInVzZSBzdHJpY3RcIjtcbiAgICBmdW5jdGlvbiBnZXRBbGwoY29udGV4dCwgdGFnKSB7XG4gICAgICAgIC8vIFN1cHBvcnQ6IElFIDw9OSAtIDExIG9ubHlcbiAgICAgICAgLy8gVXNlIHR5cGVvZiB0byBhdm9pZCB6ZXJvLWFyZ3VtZW50IG1ldGhvZCBpbnZvY2F0aW9uIG9uIGhvc3Qgb2JqZWN0cyAoIzE1MTUxKVxuICAgICAgICB2YXIgcmV0O1xuICAgICAgICBpZiAodHlwZW9mIGNvbnRleHQuZ2V0RWxlbWVudHNCeVRhZ05hbWUgIT09IFwidW5kZWZpbmVkXCIpIHtcbiAgICAgICAgICAgIHJldCA9IGNvbnRleHQuZ2V0RWxlbWVudHNCeVRhZ05hbWUodGFnIHx8IFwiKlwiKTtcbiAgICAgICAgfVxuICAgICAgICBlbHNlIGlmICh0eXBlb2YgY29udGV4dC5xdWVyeVNlbGVjdG9yQWxsICE9PSBcInVuZGVmaW5lZFwiKSB7XG4gICAgICAgICAgICByZXQgPSBjb250ZXh0LnF1ZXJ5U2VsZWN0b3JBbGwodGFnIHx8IFwiKlwiKTtcbiAgICAgICAgfVxuICAgICAgICBlbHNlIHtcbiAgICAgICAgICAgIHJldCA9IFtdO1xuICAgICAgICB9XG4gICAgICAgIGlmICh0YWcgPT09IHVuZGVmaW5lZCB8fCB0YWcgJiYgbm9kZU5hbWUoY29udGV4dCwgdGFnKSkge1xuICAgICAgICAgICAgcmV0dXJuIGpRdWVyeS5tZXJnZShbY29udGV4dF0sIHJldCk7XG4gICAgICAgIH1cbiAgICAgICAgcmV0dXJuIHJldDtcbiAgICB9XG4gICAgcmV0dXJuIGdldEFsbDtcbn0pO1xuLy8jIHNvdXJjZU1hcHBpbmdVUkw9ZGF0YTphcHBsaWNhdGlvbi9qc29uO2NoYXJzZXQ9dXRmODtiYXNlNjQsZXlKMlpYSnphVzl1SWpvekxDSnpiM1Z5WTJWeklqcGJJbmQzZHk5d1lXTnJjeTlxY1hWbGNua3ZjM0pqTDIxaGJtbHdkV3hoZEdsdmJpOW5aWFJCYkd3dWFuTWlYU3dpYm1GdFpYTWlPbHRkTENKdFlYQndhVzVuY3lJNklrRkJRVUVzVFVGQlRTeERRVUZGTzBsQlExQXNVMEZCVXp0SlFVTlVMR3RDUVVGclFqdERRVU5zUWl4RlFVRkZMRlZCUVZVc1RVRkJUU3hGUVVGRkxGRkJRVkU3U1VGRk4wSXNXVUZCV1N4RFFVRkRPMGxCUldJc1owSkJRV2xDTEU5QlFVOHNSVUZCUlN4SFFVRkhPMUZCUlRWQ0xEUkNRVUUwUWp0UlFVTTFRaXdyUlVGQkswVTdVVUZETDBVc1NVRkJTU3hIUVVGSExFTkJRVU03VVVGRlVpeEZRVUZGTEVOQlFVTXNRMEZCUlN4UFFVRlBMRTlCUVU4c1EwRkJReXh2UWtGQmIwSXNTMEZCU3l4WFFVRlpMRU5CUVVNc1EwRkJReXhEUVVGRE8xbEJRek5FTEVkQlFVY3NSMEZCUnl4UFFVRlBMRU5CUVVNc2IwSkJRVzlDTEVOQlFVVXNSMEZCUnl4SlFVRkpMRWRCUVVjc1EwRkJSU3hEUVVGRE8xRkJSV3hFTEVOQlFVTTdVVUZCUXl4SlFVRkpMRU5CUVVNc1JVRkJSU3hEUVVGRExFTkJRVVVzVDBGQlR5eFBRVUZQTEVOQlFVTXNaMEpCUVdkQ0xFdEJRVXNzVjBGQldTeERRVUZETEVOQlFVTXNRMEZCUXp0WlFVTTVSQ3hIUVVGSExFZEJRVWNzVDBGQlR5eERRVUZETEdkQ1FVRm5RaXhEUVVGRkxFZEJRVWNzU1VGQlNTeEhRVUZITEVOQlFVVXNRMEZCUXp0UlFVVTVReXhEUVVGRE8xRkJRVU1zU1VGQlNTeERRVUZETEVOQlFVTTdXVUZEVUN4SFFVRkhMRWRCUVVjc1JVRkJSU3hEUVVGRE8xRkJRMVlzUTBGQlF6dFJRVVZFTEVWQlFVVXNRMEZCUXl4RFFVRkZMRWRCUVVjc1MwRkJTeXhUUVVGVExFbEJRVWtzUjBGQlJ5eEpRVUZKTEZGQlFWRXNRMEZCUlN4UFFVRlBMRVZCUVVVc1IwRkJSeXhEUVVGSExFTkJRVU1zUTBGQlF5eERRVUZETzFsQlF6VkVMRTFCUVUwc1EwRkJReXhOUVVGTkxFTkJRVU1zUzBGQlN5eERRVUZGTEVOQlFVVXNUMEZCVHl4RFFVRkZMRVZCUVVVc1IwRkJSeXhEUVVGRkxFTkJRVU03VVVGRGVrTXNRMEZCUXp0UlFVVkVMRTFCUVUwc1EwRkJReXhIUVVGSExFTkJRVU03U1VGRFdpeERRVUZETzBsQlJVUXNUVUZCVFN4RFFVRkRMRTFCUVUwc1EwRkJRenRCUVVOa0xFTkJRVU1zUTBGQlJTeERRVUZESWl3aVptbHNaU0k2SW5kM2R5OXdZV05yY3k5cWNYVmxjbmt2YzNKakwyMWhibWx3ZFd4aGRHbHZiaTluWlhSQmJHd3Vhbk1pTENKemIzVnlZMlZ6UTI5dWRHVnVkQ0k2V3lKa1pXWnBibVVvSUZ0Y2JseDBYQ0l1TGk5amIzSmxYQ0lzWEc1Y2RGd2lMaTR2WTI5eVpTOXViMlJsVG1GdFpWd2lYRzVkTENCbWRXNWpkR2x2YmlnZ2FsRjFaWEo1TENCdWIyUmxUbUZ0WlNBcElIdGNibHh1WENKMWMyVWdjM1J5YVdOMFhDSTdYRzVjYm1aMWJtTjBhVzl1SUdkbGRFRnNiQ2dnWTI5dWRHVjRkQ3dnZEdGbklDa2dlMXh1WEc1Y2RDOHZJRk4xY0hCdmNuUTZJRWxGSUR3OU9TQXRJREV4SUc5dWJIbGNibHgwTHk4Z1ZYTmxJSFI1Y0dWdlppQjBieUJoZG05cFpDQjZaWEp2TFdGeVozVnRaVzUwSUcxbGRHaHZaQ0JwYm5adlkyRjBhVzl1SUc5dUlHaHZjM1FnYjJKcVpXTjBjeUFvSXpFMU1UVXhLVnh1WEhSMllYSWdjbVYwTzF4dVhHNWNkR2xtSUNnZ2RIbHdaVzltSUdOdmJuUmxlSFF1WjJWMFJXeGxiV1Z1ZEhOQ2VWUmhaMDVoYldVZ0lUMDlJRndpZFc1a1pXWnBibVZrWENJZ0tTQjdYRzVjZEZ4MGNtVjBJRDBnWTI5dWRHVjRkQzVuWlhSRmJHVnRaVzUwYzBKNVZHRm5UbUZ0WlNnZ2RHRm5JSHg4SUZ3aUtsd2lJQ2s3WEc1Y2JseDBmU0JsYkhObElHbG1JQ2dnZEhsd1pXOW1JR052Ym5SbGVIUXVjWFZsY25sVFpXeGxZM1J2Y2tGc2JDQWhQVDBnWENKMWJtUmxabWx1WldSY0lpQXBJSHRjYmx4MFhIUnlaWFFnUFNCamIyNTBaWGgwTG5GMVpYSjVVMlZzWldOMGIzSkJiR3dvSUhSaFp5QjhmQ0JjSWlwY0lpQXBPMXh1WEc1Y2RIMGdaV3h6WlNCN1hHNWNkRngwY21WMElEMGdXMTA3WEc1Y2RIMWNibHh1WEhScFppQW9JSFJoWnlBOVBUMGdkVzVrWldacGJtVmtJSHg4SUhSaFp5QW1KaUJ1YjJSbFRtRnRaU2dnWTI5dWRHVjRkQ3dnZEdGbklDa2dLU0I3WEc1Y2RGeDBjbVYwZFhKdUlHcFJkV1Z5ZVM1dFpYSm5aU2dnV3lCamIyNTBaWGgwSUYwc0lISmxkQ0FwTzF4dVhIUjlYRzVjYmx4MGNtVjBkWEp1SUhKbGREdGNibjFjYmx4dWNtVjBkWEp1SUdkbGRFRnNiRHRjYm4wZ0tUdGNiaUpkZlE9PVxuXG4vLyMgc291cmNlTWFwcGluZ1VSTD1kYXRhOmFwcGxpY2F0aW9uL2pzb247Y2hhcnNldD11dGY4O2Jhc2U2NCxleUoyWlhKemFXOXVJam96TENKemIzVnlZMlZ6SWpwYkluZDNkeTl3WVdOcmN5OXFjWFZsY25rdmMzSmpMMjFoYm1sd2RXeGhkR2x2Ymk5blpYUkJiR3d1YW5NaVhTd2libUZ0WlhNaU9sdGRMQ0p0WVhCd2FXNW5jeUk2SWtGQlFVRXNUVUZCVFN4RFFVRkRPMGxCUTBnc1UwRkJVenRKUVVOVUxHdENRVUZyUWp0RFFVTnlRaXhGUVVGRkxGVkJRVlVzVFVGQlRTeEZRVUZGTEZGQlFWRTdTVUZEZWtJc1dVRkJXU3hEUVVGRE8wbEJRMklzWjBKQlFXZENMRTlCUVU4c1JVRkJSU3hIUVVGSE8xRkJRM2hDTERSQ1FVRTBRanRSUVVNMVFpd3JSVUZCSzBVN1VVRkRMMFVzU1VGQlNTeEhRVUZITEVOQlFVTTdVVUZEVWl4RlFVRkZMRU5CUVVNc1EwRkJReXhQUVVGUExFOUJRVThzUTBGQlF5eHZRa0ZCYjBJc1MwRkJTeXhYUVVGWExFTkJRVU1zUTBGQlF5eERRVUZETzFsQlEzUkVMRWRCUVVjc1IwRkJSeXhQUVVGUExFTkJRVU1zYjBKQlFXOUNMRU5CUVVNc1IwRkJSeXhKUVVGSkxFZEJRVWNzUTBGQlF5eERRVUZETzFGQlEyNUVMRU5CUVVNN1VVRkRSQ3hKUVVGSkxFTkJRVU1zUlVGQlJTeERRVUZETEVOQlFVTXNUMEZCVHl4UFFVRlBMRU5CUVVNc1owSkJRV2RDTEV0QlFVc3NWMEZCVnl4RFFVRkRMRU5CUVVNc1EwRkJRenRaUVVOMlJDeEhRVUZITEVkQlFVY3NUMEZCVHl4RFFVRkRMR2RDUVVGblFpeERRVUZETEVkQlFVY3NTVUZCU1N4SFFVRkhMRU5CUVVNc1EwRkJRenRSUVVNdlF5eERRVUZETzFGQlEwUXNTVUZCU1N4RFFVRkRMRU5CUVVNN1dVRkRSaXhIUVVGSExFZEJRVWNzUlVGQlJTeERRVUZETzFGQlEySXNRMEZCUXp0UlFVTkVMRVZCUVVVc1EwRkJReXhEUVVGRExFZEJRVWNzUzBGQlN5eFRRVUZUTEVsQlFVa3NSMEZCUnl4SlFVRkpMRkZCUVZFc1EwRkJReXhQUVVGUExFVkJRVVVzUjBGQlJ5eERRVUZETEVOQlFVTXNRMEZCUXl4RFFVRkRPMWxCUTNKRUxFMUJRVTBzUTBGQlF5eE5RVUZOTEVOQlFVTXNTMEZCU3l4RFFVRkRMRU5CUVVNc1QwRkJUeXhEUVVGRExFVkJRVVVzUjBGQlJ5eERRVUZETEVOQlFVTTdVVUZEZUVNc1EwRkJRenRSUVVORUxFMUJRVTBzUTBGQlF5eEhRVUZITEVOQlFVTTdTVUZEWml4RFFVRkRPMGxCUTBRc1RVRkJUU3hEUVVGRExFMUJRVTBzUTBGQlF6dEJRVU5zUWl4RFFVRkRMRU5CUVVNc1EwRkJRenRCUVVWSUxIVjBSVUZCZFhSRklpd2labWxzWlNJNkluZDNkeTl3WVdOcmN5OXFjWFZsY25rdmMzSmpMMjFoYm1sd2RXeGhkR2x2Ymk5blpYUkJiR3d1YW5NaUxDSnpiM1Z5WTJWelEyOXVkR1Z1ZENJNld5SmtaV1pwYm1Vb1cxeHVJQ0FnSUZ3aUxpNHZZMjl5WlZ3aUxGeHVJQ0FnSUZ3aUxpNHZZMjl5WlM5dWIyUmxUbUZ0WlZ3aVhHNWRMQ0JtZFc1amRHbHZiaUFvYWxGMVpYSjVMQ0J1YjJSbFRtRnRaU2tnZTF4dUlDQWdJRndpZFhObElITjBjbWxqZEZ3aU8xeHVJQ0FnSUdaMWJtTjBhVzl1SUdkbGRFRnNiQ2hqYjI1MFpYaDBMQ0IwWVdjcElIdGNiaUFnSUNBZ0lDQWdMeThnVTNWd2NHOXlkRG9nU1VVZ1BEMDVJQzBnTVRFZ2IyNXNlVnh1SUNBZ0lDQWdJQ0F2THlCVmMyVWdkSGx3Wlc5bUlIUnZJR0YyYjJsa0lIcGxjbTh0WVhKbmRXMWxiblFnYldWMGFHOWtJR2x1ZG05allYUnBiMjRnYjI0Z2FHOXpkQ0J2WW1wbFkzUnpJQ2dqTVRVeE5URXBYRzRnSUNBZ0lDQWdJSFpoY2lCeVpYUTdYRzRnSUNBZ0lDQWdJR2xtSUNoMGVYQmxiMllnWTI5dWRHVjRkQzVuWlhSRmJHVnRaVzUwYzBKNVZHRm5UbUZ0WlNBaFBUMGdYQ0oxYm1SbFptbHVaV1JjSWlrZ2UxeHVJQ0FnSUNBZ0lDQWdJQ0FnY21WMElEMGdZMjl1ZEdWNGRDNW5aWFJGYkdWdFpXNTBjMEo1VkdGblRtRnRaU2gwWVdjZ2ZId2dYQ0lxWENJcE8xeHVJQ0FnSUNBZ0lDQjlYRzRnSUNBZ0lDQWdJR1ZzYzJVZ2FXWWdLSFI1Y0dWdlppQmpiMjUwWlhoMExuRjFaWEo1VTJWc1pXTjBiM0pCYkd3Z0lUMDlJRndpZFc1a1pXWnBibVZrWENJcElIdGNiaUFnSUNBZ0lDQWdJQ0FnSUhKbGRDQTlJR052Ym5SbGVIUXVjWFZsY25sVFpXeGxZM1J2Y2tGc2JDaDBZV2NnZkh3Z1hDSXFYQ0lwTzF4dUlDQWdJQ0FnSUNCOVhHNGdJQ0FnSUNBZ0lHVnNjMlVnZTF4dUlDQWdJQ0FnSUNBZ0lDQWdjbVYwSUQwZ1cxMDdYRzRnSUNBZ0lDQWdJSDFjYmlBZ0lDQWdJQ0FnYVdZZ0tIUmhaeUE5UFQwZ2RXNWtaV1pwYm1Wa0lIeDhJSFJoWnlBbUppQnViMlJsVG1GdFpTaGpiMjUwWlhoMExDQjBZV2NwS1NCN1hHNGdJQ0FnSUNBZ0lDQWdJQ0J5WlhSMWNtNGdhbEYxWlhKNUxtMWxjbWRsS0Z0amIyNTBaWGgwWFN3Z2NtVjBLVHRjYmlBZ0lDQWdJQ0FnZlZ4dUlDQWdJQ0FnSUNCeVpYUjFjbTRnY21WME8xeHVJQ0FnSUgxY2JpQWdJQ0J5WlhSMWNtNGdaMlYwUVd4c08xeHVmU2s3WEc1Y2JpOHZJeUJ6YjNWeVkyVk5ZWEJ3YVc1blZWSk1QV1JoZEdFNllYQndiR2xqWVhScGIyNHZhbk52Ymp0amFHRnljMlYwUFhWMFpqZzdZbUZ6WlRZMExHVjVTakphV0VwNllWYzVkVWxxYjNwTVEwcDZZak5XZVZreVZucEphbkJpU1c1a00yUjVPWGRaVjA1eVkzazVjV05ZVm14amJtdDJZek5LYWt3eU1XaGliV3gzWkZkNGFHUkhiSFppYVRsdVdsaFNRbUpIZDNWaGJrMXBXRk4zYVdKdFJuUmFXRTFwVDJ4MFpFeERTblJaV0VKM1lWYzFibU41U1RaSmEwWkNVVlZGYzFSVlJrSlVVM2hFVVZWR1JrOHdiRUpSTVVGelZUQkdRbFY2ZEVwUlZVNVZURWQwUTFGVlJuSlJhblJFVVZWT2MxRnBlRVpSVlVaR1RFWldRbEZXVlhOVVZVWkNWRk40UmxGVlJrWk1Sa1pDVVZaRk4xTlZSa1pPTUVselYxVkdRbGRUZUVSUlZVWkVUekJzUWxKWFNYTmFNRXBDVVZkc1EweEZPVUpSVlRoelVsVkdRbEpUZUVoUlZVWklUekZHUWxKVVZrTk1SRkpEVVZWRk1GRnFkRkpSVlUweFVXbDNjbEpWUmtKTE1GVTNWVlZHUkV3d1ZYTlRWVVpDVTFONFNGRlZSa2hNUlU1Q1VWVk5OMVZWUmtaVmFYaEdVVlZHUmt4RlRrSlJWVTF6VVRCR1FsSlRlRkJSVlVaUVRFVTVRbEZWT0hOUk1FWkNVWGw0ZGxGclJrSmlNRWx6VXpCR1FsTjVlRmhSVlVaYVRFVk9RbEZWVFhOUk1FWkNVWGw0UkZGVlJrUlBNV3hDVVhwT1JVeEZaRUpSVldOelVqQkdRbEo1ZUZCUlZVWlFURVZPUWxGVlRYTmlNRXBDVVZjNVEweEZUa0pSVlZWelVqQkdRbEo1ZUVwUlZVWktURVZrUWxGVlkzTlJNRVpDVWxONFJGRlZSa1JQTVVaQ1VsZDRSVXhGVGtKUlZVMDNWVlZHUWxGNWVFcFJWVVpLVEVWT1FsRlZUWE5TVlVaQ1VsTjRSRkZWUmtSTVJVNUNVVlZWYzFRd1JrSlVlWGhRVVZWR1VFeEZUa0pSVlUxeldqQktRbEZYWkVOTVJYUkNVVlZ6YzFZd1JrSlhVM2hFVVZWR1JFeEZUa0pSVlUxelVUQkdRbEY2ZEZwUlZVMDFVa040U0ZGVlJraE1SV1JDVVZWamMxUXdSa0pVZVhoRVVWVkdSRXhIWkVOUlZVWnVVV2w0UkZGVlJrWk1SV1JDVVZWamMxTlZSa0pUVTNoSVVWVkdTRXhGVGtKUlZWVnpVVEJHUWxGNmRGSlJWVlUxVVhsNFJGRlZSa1JQTVVaQ1VWVk5jMU5WUmtKVFUzaEVVVlZHUkV4RlRrSlJWVTAzVjFWR1JGVkRlRWhSVlVaSVRFVmtRbEZWWTNOU1ZVWkNVbE40UkZGVlJrUlBNVVpDVVRGWmMxRXdSa0pSZW5SU1VWVldSVXhGVmtKUlZWVnpVVEJHUWxGNWVFUlJWVVpHVEVWa1FsRlZZM05UTUVaQ1UzbDRWRkZWUmxSTVJXeENVVlZyYzFJd1JrSlNlWGhLVVZWR1NreEdSa0pSVmtWelVUQkdRbEpUZUZCUlZVWlFURVZXUWxGVlZYTlNNRVpDVW5sNFJGRlZSa2hNUlU1Q1VWVk5jMUV3UmtKUmVYaEVVVlZHUkU4eGJFSlJlbFpGVEVVeFFsRlZNSE5STUVaQ1VYbDRUbEZWUms1TVJVNUNVVlZOYzFNd1JrSlRlWGhFVVZWR1JreEZUa0pSVlZWelZEQkdRbFI1ZUVSUlZVWkdURVZXUWxGVlZYTlNNRVpDVW5sNFJGRlZSa1pNUlU1Q1VWVk5OMVZWUmtSbGEwMXpVVEJHUWxGNmRGSlJWVlpGVEVVeFFsRlZNSE5STUVaQ1VYbDRTRkZWUmtoTVJVNUNVVlZOTjFOVlJrUlhhWGhFVVZWR1JFOHdiRUpTVlZGelZGVkdRbFJUZUVSUlZVWkVURVV4UWxGVk1ITlJNRVpDVVhwMFFsRlZUbXRNUlU1Q1VWVk5jMUV3UmtKU1UzaEVVVlZHUkVscGQybGFiV3h6V2xOSk5rbHVaRE5rZVRsM1dWZE9jbU41T1hGaldGWnNZMjVyZG1NelNtcE1NakZvWW0xc2QyUlhlR2hrUjJ4MlltazVibHBZVWtKaVIzZDFZVzVOYVV4RFNucGlNMVo1V1RKV2VsRXlPWFZrUjFaMVpFTkpObGQ1U210YVYxcHdZbTFWYjBsR2RHTmliSGd3V0VOSmRVeHBPV3BpTTBwc1dFTkpjMWhITldOa1JuZHBUR2swZGxreU9YbGFVemwxWWpKU2JGUnRSblJhVm5kcFdFYzFaRXhEUW0xa1Z6VnFaRWRzZG1KcFoyZGhiRVl4V2xoS05VeERRblZpTWxKc1ZHMUdkRnBUUVhCSlNIUmpZbXg0ZFZoRFNqRmpNbFZuWXpOU2VXRlhUakJZUTBrM1dFYzFZMkp0V2pGaWJVNHdZVmM1ZFVsSFpHeGtSVVp6WWtObloxa3lPWFZrUjFZMFpFTjNaMlJIUm01SlEydG5aVEY0ZFZoSE5XTmtRemgyU1VaT01XTklRblpqYmxFMlNVVnNSa2xFZHpsUFUwRjBTVVJGZUVsSE9YVmlTR3hqWW14NE1FeDVPR2RXV0U1c1NVaFNOV05IVm5aYWFVSXdZbmxDYUdSdE9YQmFRMEkyV2xoS2RreFhSbmxhTTFaMFdsYzFNRWxITVd4a1IyaDJXa05DY0dKdVduWlpNa1l3WVZjNWRVbEhPWFZKUjJoMll6TlJaMkl5U25GYVYwNHdZM2xCYjBsNlJURk5WRlY0UzFaNGRWaElVakpaV0VsblkyMVdNRTh4ZUhWWVJ6VmpaRWRzYlVsRFoyZGtTR3gzV2xjNWJVbEhUblppYmxKc1pVaFJkVm95VmpCU1YzaHNZbGRXZFdSSVRrTmxWbEpvV2pBMWFHSlhWV2RKVkRBNVNVWjNhV1JYTld0YVYxcHdZbTFXYTFoRFNXZExVMEkzV0VjMVkyUkdlREJqYlZZd1NVUXdaMWt5T1hWa1IxWTBaRU0xYmxwWVVrWmlSMVowV2xjMU1HTXdTalZXUjBadVZHMUdkRnBUWjJka1IwWnVTVWg0T0VsR2QybExiSGRwU1VOck4xaEhOV05pYkhnd1psTkNiR0pJVG14SlIyeHRTVU5uWjJSSWJIZGFWemx0U1VkT2RtSnVVbXhsU0ZGMVkxaFdiR051YkZSYVYzaHNXVE5TZG1OclJuTmlRMEZvVUZRd1oxaERTakZpYlZKc1dtMXNkVnBYVW1OSmFVRndTVWgwWTJKc2VEQllTRko1V2xoUloxQlRRbXBpTWpVd1dsaG9NRXh1UmpGYVdFbzFWVEpXYzFwWFRqQmlNMHBDWWtkM2IwbElVbWhhZVVJNFprTkNZMGxwY0dOSmFVRndUekY0ZFZoSE5XTmtTREJuV2xkNGVscFRRamRZUnpWalpFWjRNR050VmpCSlJEQm5WekV3TjFoSE5XTmtTREZqWW14NGRWaElVbkJhYVVGdlNVaFNhRnA1UVRsUVZEQm5aRmMxYTFwWFduQmliVlpyU1VoNE9FbElVbWhhZVVGdFNtbENkV0l5VW14VWJVWjBXbE5uWjFreU9YVmtSMVkwWkVOM1oyUkhSbTVKUTJ0blMxTkNOMWhITldOa1JuZ3dZMjFXTUdSWVNuVkpSM0JTWkZkV2VXVlROWFJhV0VwdVdsTm5aMWQ1UW1waU1qVXdXbGhvTUVsR01ITkpTRXBzWkVOQmNFOHhlSFZZU0ZJNVdFYzFZMkpzZURCamJWWXdaRmhLZFVsSVNteGtSSFJqWW00eFkySnNlSFZqYlZZd1pGaEtkVWxIWkd4a1JVWnpZa1IwWTJKdU1HZExWSFJqWW1sS1pHWlJQVDFjYmlKZGZRPT1cbiJdfQ==