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
	$("#error_student_record").hide();
	$("#exam_list").hide();
	
	$("#error_exam_source").hide();
	$("#error_old_exam").hide();
	$("#success_exam").hide();

	$("#error_add_group").hide();
	$("#success_add_group").hide();

	$("#error_add_student").hide();
	$("#success_add_student").hide();

	$("#take_form").hide();

	$scope.queryAllGroups = function () {

		appFactory.queryAllGroups(function (data) {
			var array = [];
			for (var i = 0; i < data.length; i++) {
				data[i].Record.Key = data[i].Key;
				array.push(data[i].Record);
			}
			array.sort(function(a, b) {
			    return a.groupName.localeCompare(b.groupName);
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

    $scope.addStudent = function () {

		appFactory.addStudent($scope.student, function (data) {

			if (data == "Could not locate unpassed test") {
				$("#error_add_student").show();
				$("#success_add_student").hide();
			} else {
				$("#error_add_student").hide();
				$("#success_add_student").show();
			}

			$scope.exam_result = data;
		});
	}

	$scope.queryAllStudents = function () {

		appFactory.queryAllStudents(function (data) {
			var array = [];
			for (var i = 0; i < data.length; i++) {
				data[i].Record.Key = data[i].Key;
				array.push(data[i].Record);
			}
			array.sort(function (a, b) {
				return a.groupName.localeCompare(b.groupName);
			});
			$scope.all_students = array;
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

	$scope.getStudentRecord = function () {
		
		var id = $scope.id;

		appFactory.getStudentRecord(id, function(data){

			$scope.student_record = data;

			if ($scope.student_record == "Student record not found"){
				console.log()
				$("#error_student_record").show();
				$("#student_record").hide();
				$("#student_record2").hide();
				
			} else{
				$("#error_student_record").hide();
				$("#student_record").show();
				$("#student_record2").show();
			}
		});
	}

	$scope.prepareForExam = function () {

		var exam = $scope.exam;

		appFactory.prepareForExam(exam, function (data) {

			if (data == "No group/course found") {
				console.log()
				$("#error_prepare_exam").show();
				$("#exam_list").hide();
			
			} else {
				$("#error_prepare_exam").hide();
				$("#exam_list").show();
			}

			var array = [];
			for (var i = 0; i < data.length; i++) {
				//parseInt(data[i].Key);
				data[i].Record.Key = data[i].Key;
				array.push(data[i].Record);
			}
			array.sort(function (a, b) {
				return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.exam_list = array;
		});
	}


	$scope.beforeTakeTheTest = function (exam) {
		        
          if (exam.rate != "") {
			$("#takeTheTestId").hide();	 
			$("#take_form").hide(); 
		  } else {
			$("#takeTheTestId").show();	
			$("#take_form").show();
			$("#success_exam").show();
		  }
		  $scope.examcase = exam;
	}

	$scope.takeTheTest = function () {

		var examcase = $scope.examcase;

		appFactory.takeTheTest(examcase, function (data) {

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

		var newGroup =  data.groupName + "-" + data.description;

		$http.get('/add_group/' + newGroup).success(function (output) {
			callback(output)
		});
	}


	factory.addStudent = function (data, callback) {

		var student = data.studentId + "-" + data.studentName + "-" + data.groupName + "-" + data.description;

		$http.get('/add_student/' + student).success(function (output) {
			callback(output)
		});
	}

	factory.queryAllStudents = function (callback) {

		$http.get('/get_all_students/').success(function (output) {
			callback(output)
		});
	}
	
	factory.createTestForGroup = function (generator, callback) {

		var generator = generator.groupName + "-" + generator.courseName + "-" + generator.teacherName;

		$http.get('/create_test_group/' + generator).success(function (output) {
			callback(output)
		});
	}

	factory.queryTestById = function (id, callback) {
		$http.get('/get_test_id/' + id).success(function (output) {
			callback(output)
		});
	}

	factory.getStudentRecord = function (id, callback) {
		$http.get('/get_student_record/' + id).success(function (output) {
			callback(output)
		});
	}

	factory.prepareForExam = function (exam, callback) {

		var params = exam.group + "-" + exam.course;

		$http.get('/prepare_exam/' + params).success(function (output) {
			callback(output)
		});
	}

	factory.takeTheTest = function (input, callback) {

		var params = input.studentId + "-" + input.course + "-" + input.rate;

		$http.get('/take_test/' + params).success(function (output) {
			callback(output)
		});
	}

	return factory;
});
