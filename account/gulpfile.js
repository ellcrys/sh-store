'use strict';
 
var gulp = require('gulp');
var sass = require('gulp-sass');
var rename = require('gulp-rename');
var nodemon = require('gulp-nodemon')
 
// Sass
gulp.task('sass', function () {
  return gulp.src('./assets/styles/importer.scss')
    .pipe(sass().on('error', sass.logError))
    .pipe(rename("all.css"))
    .pipe(gulp.dest('./assets/styles'));
});
 
// Sass watch 
gulp.task('sass:watch', function () {
  gulp.watch('./assets/styles/**/*.scss', ['sass']);
});

// Sass serve
gulp.task('serve', function name() {
    nodemon({
        script: 'app.js',
        ext: 'js',
        env: { 'NODE_ENV': 'development' }
    })
})

gulp.task('develop', ["sass:watch", "serve"])