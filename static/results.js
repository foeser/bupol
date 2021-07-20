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
        console.log(data)
        window.appData = [data];
        renderPersonsData();
    }).fail(function (data) {
        alert("Error: Problem with getting approver information.");
    });

}

function renderPersonsData() {

    var BSControl = function (config) {
        jsGrid.ControlField.call(this, config);
    };

    BSControl.prototype = new jsGrid.ControlField({
        _createUpdateButton: function () {
            var grid = this._grid;
            return $("<button class=\"btn btn-sm btn-light m-0 ml-1 p-1\" title=\"" + this.updateButtonTooltip + "\">").append("<i class=\"fas fa-check bs-grid-button text-success m-0 p-0\">").click(function (e) {
                console.error("onUpdate")
                grid.updateItem();
                e.stopPropagation();
            });
        }
    });

    jsGrid.fields.bscontrol = BSControl;

    var gridFields = [];
    var gridWidth = "1700px";

    gridFields.push({ name: GridfieldName1, title: GridfieldTitle1, type: "text", width: 150, validate: { validator: "required", message: "Device name is a required field." },
        editing: false,
        align: "center"

    });
    gridFields.push({ name: GridfieldName2, title: GridfieldTitle2, type: "text", width: 150, validate: { validator: "required", message: "Device name is a required field." },
        // show empty input field while editing and allow enter key for validation
        editTemplate: function (value, item) {
            var $result = jsGrid.fields.text.prototype.editTemplate.call(this);
            $result.on("keydown", function(e) {
                if(e.which === 13) {
                    $("#GridPersons").jsGrid("updateItem");
                    return false;
                }
            });
            return $result;
        },
        align: "center"
    });
    gridFields.push({ name: "editable", title: "Editable", type: "text", width: 150, editing: false,visible: false });
    gridFields.push({
        name: "command", type: "bscontrol", width: 125, modeSwitchButton: false, editing: false, inserting: false,editButton: false, deleteButton: false,
        itemTemplate: function (value, item) {
            if(item.editable === "false") {
                return ""
            }
            return $("<button>").addClass("btn btn-primary btn-sm")
                .attr({ type: "button", title: "Click in row to start!" })
                .html("Click in row to start!")
        }
    });

    $("#GridPersons").jsGrid({
        height: "auto",
        width:  gridWidth,
        updateOnResize: true,
        editing: true,
        inserting: false,
        sorting: false,
        align: "center",
        data: appData,
        fields: gridFields,
        controller: {
            // check if entered values are correct
            updateItem: function(item) {
                return $.ajax({
                    type: "GET",
                    url: vDir + "/query/" + Exercise + "/" + item[GridfieldName1] + '/' + item[GridfieldName2],
                    contentType: "application/json; charset=utf-8",
                    dataType: "json",
                    success: function (data) {
                        console.log("Match")
                        data.message = "match"
                        //$.showSnackBar(data);
                        return data

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
            }
        },
        // get next item/row
        onItemUpdated: function(args) {
            $.ajax({
                type: "GET",
                url:vDir + "/data/getResults/" + Exercise + "/" + (parseInt(args.itemIndex) + 1),
                contentType: "text/plain",
                dataType: "json",
                success: function (data) {
                    //$.showSnackBar(data);
                    console.log(data)
                    $("#GridPersons").jsGrid("insertItem", data).done(function() {
                        console.log("insertion completed");
                    });
                },
                error: function (data, err) {
                    //var json = $.parseJSON(data.responseJSON);
                    //$.showSnackBar(json);
                    console.log(data)
                    //console.error("Message: " + json.message);
                    if(data.responseJSON.message === "Index out of range.") {
                        $("#GridPersons").jsGrid("option", "editing", false)
                        data.message = "YOU SOLVED EVERYTHING! :)"
                        $.showSnackBar(data);
                    }
                    console.error(data.responseJSON.message);
                }
            })
            /*console.log("row " + args.itemIndex)
            $.getJSON(vDir + "/data/getResults/" + Exercise + "/" + (parseInt(args.itemIndex) + 1), function (data) {
                console.log(data)
                $("#GridPersons").jsGrid("insertItem", data).done(function() {
                    console.log("insertion completed");
                });
            }).fail(function (data) {
                $.showSnackBar(data);
                console.error("Message: " + data.message);
                $("#GridPersons").jsGrid("option", "editing", false)
                console.error(data);
            });*/
        },
    });
}