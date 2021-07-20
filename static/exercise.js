$(document).ready(function () {
    getAppData();
    startTimer();
});

function startTimer() {
    var timer = new Timer();
    timer.start({countdown: true, startValues: {seconds: parseInt(timerTime, 10)}});

    $('#countdown .values').html(timer.getTimeValues().toString());

    timer.addEventListener('secondsUpdated', function (e) {
        $('#countdown .values').html(timer.getTimeValues().toString());
    });

    timer.addEventListener('targetAchieved', function (e) {
        //$('#countdown .values').html('KABOOM!!');
        //$("#GridPersons").children().remove();
        $(".jsgrid-grid-body").children().remove();
    });
}

function getAppData() {
    $.getJSON(vDir + "/data/get/" + Exercise, function (data) {
        //console.log(data)
        window.appData = data;
        renderPersonsData();
    }).fail(function (data) {
        alert("Error: Problem with getting approver information.");
    });
}

function renderPersonsData() {

    var gridFields = [];
    var gridWidth = "700px";

    gridFields.push({ name: GridfieldName1, title: GridfieldTitle1, type: "text", width: 150, validate: { validator: "required", message: "Device name is a required field." },
        editing: false,
        align: "center"

    });
    gridFields.push({ name: GridfieldName2, title: GridfieldTitle2, type: "text", width: 150, validate: { validator: "required", message: "Device name is a required field." },
        editing: false,
        align: "center"
    });

    $("#GridPersons").jsGrid({
        height: "auto",
        width:  gridWidth,
        updateOnResize: true,
        editing: false,
        inserting: false,
        sorting: false,
        confirmDeleting: true,
        align: "center",
        data: appData,
        fields: gridFields,
    });
}

function ReloadData() {
   getAppData()
}


