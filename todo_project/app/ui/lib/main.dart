import 'package:flutter/material.dart';

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

  String getTodo() {
    
    return "";
  }
  @override
  Widget build(BuildContext context) {
    // TODO: implement build
    return Scaffold(
      backgroundColor: backgroundCol,
      body: Center(
        child: Container(
          width: MediaQuery.of(context).size.width/1.5,
          child: ListView.builder(itemBuilder: itemBuilder),
        ),
      ),
      appBar: AppBar(
        backgroundColor: appBarCol,
        title: Text("MyTodo"),
        centerTitle: true,
      ),

      floatingActionButton: FloatingActionButton(
        onPressed: (){},
        backgroundColor: primary,
        child: Icon(Icons.add),
      ),
    );
  }
}
