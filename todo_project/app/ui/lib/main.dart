// import 'dart:ffi';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

void main() {
  runApp(Todo());
}

class Todo extends StatelessWidget {
  // This widget is the root of your application.
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'MyTodo',
      debugShowCheckedModeBanner: false,
      home: HomePage(),
    );
  }
}

class HomePage extends StatefulWidget {

  @override
  _HomePageState createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {

  Color backgroundCol = new Color(0xff222831);
  Color appBarCol = new Color(0xff393e46);
  Color primary = new Color(0xff00adb5);
  Color tiles = new Color(0xff3e4149);
  Color secondary = new Color(0xffeeeeee);
  Color listText = new Color(0xffdbe2ef);
  Future<http.Response> getTodo() {
    
    Future<http.Response> resp = http.get("http://127.0.0.1:23556",
    headers: <String, String>{
      'Content-Type': 'application/json; charset=UTF-8',
    },);
    return resp;
  }
  @override
  void initState() {
    // TODO: implement initState
    super.initState();
    getTodo();
  }
  @override
  Widget build(BuildContext context) {
    // TODO: implement build
    return Scaffold(
      backgroundColor: backgroundCol,
      body: Center(
        child: Container(
          width: MediaQuery.of(context).size.width/3.5,
          child: FutureBuilder<http.Response>(
            future: getTodo(),
            builder: (context, snapshot){
              if(snapshot.hasData){
                var responseString = json.decode(snapshot.data.body);
                var incomplete = responseString['incomplete'].toString().split("\n");
                var complete = responseString['complete'].toString().split("\n");
               
                return Container(
                  margin: EdgeInsets.only(top: 25),
                  child: ListView.builder(
                    physics: BouncingScrollPhysics(),
                    itemCount: incomplete.length-1,
                    itemBuilder: (context, index){

                      final item = incomplete[index];
                      return Dismissible(
                        key: Key(item),
                        onDismissed: (direction){
                          setState(){
                            incomplete.removeAt(index);
                          }
                        },
                        child: Card(
                          margin: EdgeInsets.all(5),
                          clipBehavior: Clip.antiAlias,
                          color: tiles,
                          elevation: 3,
                          child: ListTile(
                            title: Text(incomplete[index], style: TextStyle(color: listText),),
                          ),
                          shape: StadiumBorder(
                            
                          ),
                        ),
                      );
                    },
                   ),
                  // child: AnimatedList(
                  //   initialItemCount: incomplete.length-1,
                  //   itemBuilder: (context, index, animation){
                  //     return SlideTransition(
                  //       position: animation.drive(Tween<Offset>(begin: Offset(-1, 0), end: Offset(0, 0))),
                  //       child: Card(
                  //       margin: EdgeInsets.all(5),
                  //       clipBehavior: Clip.antiAlias,
                  //       color: tiles,
                  //       elevation: 3,
                  //       child: ListTile(
                  //         contentPadding: EdgeInsets.symmetric(vertical: 7.5, horizontal: 20),
                  //         title: Text(incomplete[index], style: TextStyle(color: listText),),
                  //       ),
                  //       shape: StadiumBorder(
                          
                  //       ),
                  //     )
                  //     );
                  //   },
                  // ),
                );
              }
              else {
                return Center(child: CircularProgressIndicator());
              }
            },
            )
        ),
      ),
      appBar: AppBar(
        backgroundColor: appBarCol,
        title: Text("MyTodo"),
        centerTitle: true,
      )
    );
  }
}
