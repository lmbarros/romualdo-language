extends Node2D

var passageClass = preload("res://TestPassage.gd")

func _ready():
	$Console.text += "Initialzing Passage...\n"
	var myVarIs := "MyVar is %s\n"
	var passage = passageClass.new()
	
	$Console.text += myVarIs % passage.MyVar

	$Console.text += "Calling Passage...\n"
	passage.ChangeMyVar()
	
	$Console.text += myVarIs % passage.MyVar
