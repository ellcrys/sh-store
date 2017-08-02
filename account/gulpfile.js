'use strict';
var gulp = require('gulp');
var sass = require('gulp-sass');
var rename = require('gulp-rename');
var nodemon = require('gulp-nodemon');
var ts = require('gulp-typescript');
var sourcemaps = require("gulp-sourcemaps");
var tsProject = ts.createProject('tsconfig.json');
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
        ext: 'ts',
        tasks: ["ts"],
        ignore: ["node_modules/"],
        env: { 'NODE_ENV': 'development' }
    });
});
gulp.task('ts', function () {
    var tsResult = gulp.src(["./**/*.ts", '!./node_modules/**'], { base: '.' })
        .pipe(sourcemaps.init())
        .pipe(tsProject()).js
        .pipe(sourcemaps.write())
        .pipe(gulp.dest('.'));
});
gulp.task('ts:watch', ['ts'], function () {
    gulp.watch('**/*.ts', ['ts']);
});
gulp.task('develop', ["sass:watch", "serve"]);
