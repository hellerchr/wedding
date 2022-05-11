$(document).ready(function () {
    renderTable();
});

function renderTable() {
    $.get("/admin/rsvp", function (rsvps) {
        var hiddenFields = {"ID": true, "SessionID": true, "PersonNames": true};

        var headerFormatter = {
            "Accept": "Kommt",
            "PersonCount": "Personen",
            "Children": "Kinder",
            "Diet": "Essen",
            "Message": "Nachricht",
            "CreatedOn": "Rückmeldung am"
        };

        var valueFormatter = {
            "CreatedOn": function (date) {
                // number as double digit
                function dd(n) {
                    return n < 10 ? "0" + n : n;
                }

                var d = new Date(date);
                return dd(d.getDate()) + "." + dd(d.getMonth() + 1) + "." + (d.getFullYear()) + " " + dd(d.getHours()) + ":" + dd(d.getMinutes())
            },
            "Accept": function (value) {
                return value ? "✓" : "✗";
            }
        };

        // clear table
        var table = $("#rsvps");
        table.html("");

        // table head
        var thead = $("<thead>");
        var tr = $("<tr>");
        thead.append(tr);
        var rsvp = rsvps[0];
        for (var field in rsvp) {
            if (hiddenFields.hasOwnProperty(field)) {
                continue;
            }
            var value = field;
            if (headerFormatter.hasOwnProperty(field)) {
                value = headerFormatter[field];
            }
            tr.append('<th scope="col">' + value + '</th>')
        }
        tr.append('<th scope="col">Aktion</th>')
        table.append(thead);

        // table body
        var tbody = $("<tbody>");
        for (var i = 0; i < rsvps.length; i++) {
            var tr = $("<tr>");
            var rsvp = rsvps[i];
            for (var field in rsvp) {
                if (hiddenFields.hasOwnProperty(field)) {
                    continue;
                }
                var value = rsvp[field];
                if (valueFormatter.hasOwnProperty(field)) {
                    value = valueFormatter[field](value)
                }
                tr.append("<td>" + value + "</td>")
            }

            var del = $("<td><a href='/'>löschen</a></td>");
            del.attr("id", rsvp["ID"])
            del.click(function (e) {
                e.preventDefault();
                if (confirm("Eintrag du diesen Eintrag wirklich löschen?")) {
                    var id = $(this).attr("id");
                    return $.ajax({
                        url: "/admin/rsvp/" + id,
                        type: 'DELETE',
                        success: function () {
                            renderTable();
                        }
                    });
                }
            });
            tr.append(del);
            tbody.append(tr)
        }
        table.append(tbody);
    });
}

/*
<thead>
        <tr>
            <th scope="col">#</th>
            <th scope="col">First</th>
            <th scope="col">Last</th>
            <th scope="col">Handle</th>
        </tr>
        </thead>

<tbody>
    <tr>
      <th scope="row">1</th>
      <td>Mark</td>
      <td>Otto</td>
      <td>@mdo</td>
    </tr>
    <tr>
      <th scope="row">2</th>
      <td>Jacob</td>
      <td>Thornton</td>
      <td>@fat</td>
    </tr>
    <tr>
      <th scope="row">3</th>
      <td>Larry</td>
      <td>the Bird</td>
      <td>@twitter</td>
    </tr>
  </tbody>
 */