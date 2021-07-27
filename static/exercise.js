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
        // remove grid once countdown exceed
        $(".jsgrid-grid-body").children().remove();
    });
}

function getAppData() {
    $.getJSON(vDir + "/data/get/" + Exercise, function (data) {
        window.appData = data;
        renderPersonsData();
    }).fail(function (data) {
        alert("Error: Problem with getting information from backend.");
    });
}

function renderPersonsData() {

    var gridFields = [];

    gridFields.push({ name: GridfieldName1, title: GridfieldTitle1, type: "text", editing: false, align: "center" });
    gridFields.push({ name: GridfieldName2, title: GridfieldTitle2, type: "text", editing: false, align: "center" });

    $("#GridExercise").jsGrid({
        height: "auto",
        width: "50%",
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


