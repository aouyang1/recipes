function getTags() {
    d3.request("/tags")
        .get(function(error, data) {
            if (error) throw error;
            tags = JSON.parse(data.response);
            store.listItems = tags
            if (tags) {
                clearListItems()
                clearRecipeUpdate();
                clearListSubItems();

                renderNewTagButton();
                renderTagList();
            }
        })
}

function renderNewTagButton() {
    /*
    <div class="btn-group">
      <button class="btn btn-secondary btn-sm dropdown-toggle" type="button" data-toggle="dropdown">
        New/Link 
      </button>
      <div class="dropdown-menu">
        <form class="px-4 py-3">
          <div class="form-group">
            <input type="text" class="form-control" id="dropdownFormTagName" placeholder="tag name">
          </div>
          <button type="button" class="btn btn-primary">Create</button>
        </form>
      </div>
    </div>
    */
    list = d3.select("#list-items");

    btnGroup = list.append("div")
            .attr("class", "btn-group px-4 pb-4");

    btnGroup
        .append("button")
            .attr("class", "btn btn-secondary btn-sm dropdown-toggle")
            .attr("type", "button")
            .attr("data-bs-toggle", "dropdown")
            .text("New Tag");

    btnForm = btnGroup
        .append("div")
            .attr("class", "dropdown-menu")
        .append("form")
            .attr("class", "px-1");

    btnForm
        .append("div")
            .attr("class", "form-group py-1")
        .append("input")
            .attr("type", "text")
            .attr("class", "form-control")
            .attr("id", "dropdownFormTagName")
            .attr("placeholder", "tag name");

    btnForm
        .append("button")
            .attr("type", "button")
            .attr("class", "btn btn-primary py-2")
            .text("Create")
        .on("click", (_, d) => {
            // if successful creation of recipe to event append to recipes
            // variable and focus on the newly created one
            tag_name = d3.select("#dropdownFormTagName").property('value');
            createTag(tag_name);
        });
}

function createTag(tag_name) {
   req = {
        "name": tag_name,
    };

    d3.request("/tag")
        .post(JSON.stringify(req), function(error, data) {
            if (error) {
                console.log(error);
                return
            }
            if (data.response) {
                tag = JSON.parse(data.response); 
                store.listItems.push(tag);
                renderTagList();
            }
        });
}

function renderTagList() {
   /*
    <a href="#" class="list-group-item list-group-item-action py-3 lh-tight active ">
        <div class="col-10 mb-1 small">Some placeholder content in a paragraph below the heading and date.</div>
    </a>
    */

    list = d3.select("#list-items");
    list.selectAll("a").remove();

    console.log(store.listItems);
    list.selectAll("a")
        .data(store.listItems)
        .enter()
        .append("a")
            .attr("class", "list-group-item list-group-item-action py-3 lh-tight")
            .attr("id", d => d.name + ":" + d.variant)
            .attr("data-bs-toggle", "list")
            .on("click", (_, recipe) => {
                clearListSubItems();
                clearRecipeUpdate();
            })
        .append("d")
            .attr("class", "col-10 mb-1 small")
            .html((d) => {
                return "<div class=\"row\">"+
                "<div class=\"col-11\">"+d.name+"</div>"+
                "<div class=\"col-1\">"+(d.count > 0 ? d.count : "")+"</div>"+
                "</div>"
            });
}