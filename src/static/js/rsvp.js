$(document).ready(function () {
    // form submit handler
    $("#rsvp-form").submit(onSubmit);

    // mirror name to person1 name input
    $("#name").on("input", function () {
        $("#name_person1").attr("placeholder", $(this).val());
    });

    // only show n number of name inputs
    $("#person_count").on("input", function () {
        var val = $(this).val();
        showNameInputs(val);
    });
});

function showNameInputs(count) {
    for (var i = 1; i <= 8; i++) {
        var field = $("#name_person" + i);
        if (i <= count) {
            field.removeClass("d-none");
        } else {
            field.addClass("d-none");
        }
    }
}

function onSubmit(e) {
    // supress default form action
    e.preventDefault();

    // reset error markers
    $("#rsvp :input").removeClass("is-invalid");

    function onSuccess() {
        // close dialog
        $('#rsvp').modal('hide');
        showMessage("Rückmeldung", "Vielen Dank für Deine Rückmeldung!");
    }

    function onError(res) {
        var errors = res.responseJSON;

        // highlight fields with errors
        for (var i = 0; i < errors.length; i++) {
            var field = errors[i].field;
            var formField = $("#rsvp-form *[name=" + field + "]");
            formField.addClass("is-invalid");
        }
    }

    var formData = extractFormData($(this));

    $.post("/api/rsvp", formData)
        .fail(onError)
        .done(onSuccess);
}

function extractFormData(form) {
    var formData = {};
    form.serializeArray().forEach(function (kv) {
        formData[kv.name] = kv.value;
    });
    return formData;
}

function showMessage(title, text) {
    $("#modal-title").html(title);
    $("#modal-text").html(text);
    $('#modal').modal('show');
}