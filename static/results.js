$(document).ready(function () {

    jQuery.showSnackBar = function (data) {

        $('#snackbar').text(data.message);
        if (data.error != null) {
            $('#snackbar').addClass('alert-danger');
            $('#snackbar').removeClass('alert-success')
        } else {
            $('#snackbar').removeClass('alert-danger')
            $('#snackbar').addClass('alert-success')
        }
        $('#snackbar').show();

        // After 2 seconds, hide the Div Again
        setTimeout(function () {
            $('#snackbar').hide();
        }, 2000);
    };

    getAppData();
});

function getAppData() {
    // just get first element
    $.getJSON(vDir + "/data/getResults/" + Exercise + "/0", function (data) {
        //console.log(data)
        window.appData = [data];
        renderPersonsData();
    }).fail(function (data) {
        alert("Error: Problem with getting information from backend.");
    });

}

function renderPersonsData() {

    var gridFields = [];

    // ToDo: consider using custom validators for checking actual results: http://js-grid.com/docs/#custom-validators
    /*jsGrid.validators.foobar = {
        message: "Please enter a valid time, between 00:00 and 23:59",
        validator: function(value, item) {
            console.log("Value " + value + " for item" + item)
            return false;
        }
    }*/

    gridFields.push({ name: GridfieldName1, title: GridfieldTitle1, type: "text", editing: false, align: "center" });
    gridFields.push({ name: GridfieldName2, title: GridfieldTitle2, type: "text",validate: { validator: "required", message: "Empty input is not allowed" },
        /*validate: {
            validator: "foobar",
        },*/
        // show empty input field while editing and allow enter key for validation
        editTemplate: function (value, item) {
            var $result = jsGrid.fields.text.prototype.editTemplate.call(this);
            // enter updates row (get to next row in case of success)
            $result.on("keydown", function(e) {
                if(e.which === 13) {
                    $("#GridExercise").jsGrid("updateItem");
                    return false;
                }
            });
            // set focus to textbox when clicked in row
            setTimeout(function() {
                $result.focus();
            });
            return $result;
        },
        align: "center"
    });
    gridFields.push({ name: "editable", title: "Editable", type: "text", editing: false,visible: false });
    gridFields.push({ name: "skipped", title: "Editable", type: "text", editing: false,visible: false });
    gridFields.push({
        name: "command", type: "control", modeSwitchButton: false, editing: false, inserting: false,editButton: false, deleteButton: false,
        itemTemplate: function (value, item) {
            // don't show button/control when row is solved (disabled for editing) already
            if(item != undefined && item.editable === "false") {
                return ""
            }
            return this._createGridButton("jsgrid-edit-button", "Click to Edit this row", function(grid) {
            });
        },
        editTemplate: function (value, item) {
            // set skipped to retry to allow updating (checking for result)
            item.skipped = "retry"
            var $result = this._createGridButton("jsgrid-update-button", "Click to save this row", function(grid) {
                // ToDo: skip setting itemvalue, update value instead via updateItem call
                $("#GridExercise").jsGrid("updateItem");
            });
            return $result.add(this._createGridButton("jsgrid-cancel-button", "Click to cancel this row", function(grid) {
                $("#GridExercise").jsGrid("cancelEdit");
            }))
        },
        align: "center"
    });
    gridFields.push({
        name: "skipbutton", type: "control",  modeSwitchButton: false, editing: false, inserting: false,editButton: false, deleteButton: false,
        itemTemplate: function (value, item) {
            // don't show button/control when row is solved (disabled for editing) already or when its the last row or skipped
            var items = $("#GridExercise").jsGrid("option", "data");
            var arrayLength = items.length;
            if(item.editable === "false" || arrayLength == RowCount || item.skipped === "skipped") {
                return ""
            }

            return $("<button>").addClass("btn btn-primary btn-sm")
                .attr({ type: "button", title: "Skip" })
                .html("<i class=\"fas fa-forward\"></i>Skip")
                .on("click", function () {
                    $("#GridExercise").jsGrid("updateItem", item, { GridfieldName1: item[GridfieldName1], GridfieldName2: item[GridfieldName2], editable: item.editable, skipped: "skipped" }).done(function () {
                        $.ajax({
                            type: "GET",
                            url:vDir + "/data/getResults/" + Exercise + "/" + arrayLength,
                            contentType: "text/plain",
                            dataType: "json",
                            success: function (data) {
                                console.log(data)
                                $("#GridExercise").jsGrid("insertItem", data).done(function() {
                                    console.log("insertion completed");
                                });

                            },
                            error: function (data, err) {
                                console.log(data)
                                console.error(data.responseJSON.message);
                            }
                        })
                   })
                });
        },
        // remove botton/control when editing/inserting
        editTemplate: function (value, item) { return "" },
        insertTemplate: function () { return "" },
        align: "center"
    });

    $("#GridExercise").jsGrid({
        height: "auto",
        width:  "70%",
        shrinkToFit : true,
        updateOnResize: true,
        editing: true,
        inserting: false,
        sorting: false,
        align: "center",
        data: appData,
        fields: gridFields,
        // render solved row green
        rowClass: function(item, itemIndex) {
            console.log("renderRow")
            if(item.editable === "false"){
                console.log("renderRowHighlight")
                return 'highlight-green'
            }
        },


        controller: {
            // check if entered values are correct
            updateItem: function(item) {
                if(item.editable === "false" || item.skipped === "skipped") {
                    return
                }
                return $.ajax({
                    type: "GET",
                    url: vDir + "/query/" + Exercise + "/" + item[GridfieldName1] + '/' + item[GridfieldName2],
                    contentType: "application/json; charset=utf-8",
                    dataType: "json",
                    success: function (data) {
                        console.log("Match")
                        data.message = "match"
                        $.showSnackBar(data);

                    },
                    error: function (data, err) {
                        data.message = "Doesn't match"
                        data.error = "Doesn't match"
                        $.showSnackBar(data);
                        console.error(data);
                    }
                })
            }
        },

        // rows/items which are solved shouldn't be editable anymore
        onItemEditing: function(args) {
            if(args.item.editable === "false") {
                args.cancel = true;
                return
            }
        },

        // get next item/row
        onItemUpdated: function(args) {
            if(args.item.skipped === "skipped" || args.item.skipped === "retry") {
                return
            }
            // check if everything is solved or if we need to add more rows
            var items = $("#GridExercise").jsGrid("option", "data");
            var arrayLength = items.length;

            if(arrayLength == RowCount) {
                var solvedEverything = true
                for (var i = 0; i < arrayLength; i++) {
                    console.log(items[i])
                    if (items[i].editable === "true") {
                        solvedEverything = false
                    }
                }
                if (solvedEverything) {
                    $("#GridExercise").jsGrid("option", "editing", false)
                    var data = {message: "YOU SOLVED EVERYTHING! :)"}
                    $.showSnackBar(data);

                }
                return
            }

            // if the current updatedItem is not the last one then we skip adding more (since this is an item which was
            // skipped previously and there was already an insert afterwards)
            if(args.itemIndex != arrayLength - 1) {
                return
            }

            $.ajax({
                type: "GET",
                url: vDir + "/data/getResults/" + Exercise + "/" + (parseInt(args.itemIndex) + 1),
                contentType: "text/plain",
                dataType: "json",
                success: function (data) {
                    //$.showSnackBar(data);
                    console.log(data)
                    $("#GridExercise").jsGrid("insertItem", data).done(function () {
                        console.log("insertion completed");
                    });

                },
                error: function (data, err) {
                    console.log(data)
                    console.error(data.responseJSON.message);
                }
            })
        },
    });
}