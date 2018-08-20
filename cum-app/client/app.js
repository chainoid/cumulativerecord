// SPDX-License-Identifier: Apache-2.0

'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function ($scope, appFactory) {

	$("#success_generated").hide();
	$("#error_generated").hide();
	$("#error_query").hide();
	$("#error_sender").hide();
	$("#error_query_id").hide();
	$("#error_query_student").hide();
	$("#error_prepare_exam").hide();
	$("#error_pass_exam").hide();

	$("#error_exam_source").hide();
	$("#error_old_exam").hide();
	$("#success_exam").hide();

	$("#error_add_group").hide();
	$("#success_add_group").hide();

	$scope.queryAllGroups = function () {

		appFactory.queryAllGroups(function (data) {
			var array = [];
			for (var i = 0; i < data.length; i++) {
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function (a, b) {
				return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_groups = array;
		});
	}

	$scope.addGroup = function () {

		appFactory.addGroup($scope.newGroup, function (data) {

			if (data == "Could not locate unpassed test") {
				$("#error_add_group").show();
				$("#success_add_group").hide();
			} else {
				$("#error_add_group").hide();
				$("#success_add_group").show();
			}

			$scope.exam_result = data;
		});
	}


	$scope.queryAllTests = function () {

		appFactory.queryAllTests(function (data) {
			var array = [];
			for (var i = 0; i < data.length; i++) {
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function (a, b) {
				return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_tests = array;
		});
	}

	$scope.createTestForGroup = function () {

		appFactory.createTestForGroup($scope.generator, function (data) {
			$scope.generated_test_group = data;

			if ($scope.generated_test_group == "error_generated") {
				console.log()
				$("#error_generated").show();
			} else {
				$("#error_generated").hide();
				$("#success_generated").show();

				$scope.takeTheTest = function () {

					var progress = $scope.progress;

					appFactory.takeTheTest(progress, function (data) {

						if (data == "Could not locate unpassed test") {
							$("#error_exam_source").show();
							$("#success_exam").hide();
						} else {
							$("#error_exam_source").hide();
							$("#success_exam").show();
						}

						$scope.exam_result = data;
					});
				}
			}

		});
	}

	$scope.queryTestById = function () {

		var id = $scope.test_id;

		appFactory.queryTestById(id, function (data) {
			$scope.query_test_id = data;

			if ($scope.query_test_id == "Could not locate test") {
				console.log()
				$("#error_query_id").show();
			} else {
				$("#error_query_id").hide();
			}
		});
	}

	$scope.queryTestByStudent = function () {

		var name = $scope.student_name;

		appFactory.queryTestByStudent(name, function (data) {

			var array = [];

			if (data == "No tests for student") {
				console.log()
				$("#error_query_student").show();
			} else {
				$("#error_query_student").hide();

				for (var i = 0; i < data.length; i++) {
					parseInt(data[i].Key);
					data[i].Record.Key = parseInt(data[i].Key);
					array.push(data[i].Record);
				}
				array.sort(function (a, b) {
					return parseFloat(a.Key) - parseFloat(b.Key);
				});

			}

			$scope.student_tests = array;


		});
	}

	$scope.prepareForExam = function () {

		var exam = $scope.exam;

		appFactory.prepareForExam(exam, function (data) {

			if (data == "No group/course found") {
				console.log()
				$("#error_prepare_exam").show();
			} else {
				$("#error_prepare_exam").hide();
			}

			var array = [];
			for (var i = 0; i < data.length; i++) {
				parseInt(data[i].Key);
				data[i].Record.Key = parseInt(data[i].Key);
				array.push(data[i].Record);
			}
			array.sort(function (a, b) {
				return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.exam_list = array;
		});
	}

	$scope.takeTheTest = function () {

		var progress = $scope.progress;

		appFactory.takeTheTest(progress, function (data) {

			if (data == "Could not locate unpassed test") {
				$("#error_exam_source").show();
				$("#success_exam").hide();
			} else {
				$("#error_exam_source").hide();
				$("#success_exam").show();
			}

			$scope.exam_result = data;
		});
	}

});

// Angular Factory
app.factory('appFactory', function ($http) {

	var factory = {};

	factory.queryAllGroups = function (callback) {

		$http.get('/get_all_groups/').success(function (output) {
			callback(output)
		});
	}


	factory.addGroup = function (data, callback) {

		var newGroup = data.key + "-" + data.groupId + "-" + data.groupName + "-" + data.groupDesc;

		$http.get('/add_group/' + newGroup).success(function (output) {
			callback(output)
		});
	}


	factory.queryAllTests = function (callback) {

		$http.get('/get_all_tests/').success(function (output) {
			callback(output)
		});
	}

	factory.createTestForGroup = function (generator, callback) {

		var generator = generator.key + "-" + generator.groupId + "-" + generator.groupSize + "-" + generator.courseName + "-" + generator.teacherName + "-" + generator.deadlineTS;

		$http.get('/create_test_group/' + generator).success(function (output) {
			callback(output)
		});
	}

	factory.queryTestById = function (id, callback) {
		$http.get('/get_test_id/' + id).success(function (output) {
			callback(output)
		});
	}

	factory.queryTestByStudent = function (name, callback) {
		$http.get('/get_student_test_list/' + name).success(function (output) {
			callback(output)
		});
	}

	factory.prepareForExam = function (exam, callback) {

		var params = exam.group + "-" + exam.course;

		$http.get('/prepare_exam/' + params).success(function (output) {
			callback(output)
		});
	}

	factory.takeTheTest = function (data, callback) {

		var testCase = data.testId + "-" + data.student + "-" + data.course + "-" + data.rate;

		$http.get('/take_test/' + testCase).success(function (output) {
			callback(output)
		});
	}

	return factory;
});
