//SPDX-License-Identifier: Apache-2.0

var controller = require('./controller.js');

module.exports = function(app){

  // The cum-group area
 app.get('/get_all_groups', function(req, res){
    controller.get_all_groups(req, res);
  });
  app.get('/add_group/:newGroup', function(req, res){
    controller.add_group(req, res);
  });

  // The cum-rec area
  app.get('/add_student/:student', function(req, res){
    controller.add_student(req, res);
  });
  app.get('/get_all_students', function(req, res){
    controller.get_all_students(req, res);
  });
  app.get('/create_test_group/:generator', function(req, res){
    controller.create_test_group(req, res);
  });
  app.get('/add_student/:student', function(req, res){
    controller.add_student(req, res);
  });
  app.get('/get_test_id/:id', function(req, res){
    controller.get_test_id(req, res);
  });
  app.get('/get_student_record/:id', function(req, res){
      controller.get_student_record(req, res);
  });
  app.get('/prepare_exam/:exam', function(req, res){
      controller.prepare_exam(req, res);
  });
  app.get('/take_test/:examcase', function(req, res){
      controller.take_test(req, res);
  });

}
