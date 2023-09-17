// creates a new button for the list of items
function renderListNewButton(createFunc) {
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
            .attr("id", "dropdownFormName")
            .attr("placeholder", "name");

    btnForm
        .append("button")
            .attr("type", "button")
            .attr("class", "btn btn-primary py-2")
            .text("Create")
        .on("click", (_, d) => {
            // if successful creation of recipe to event append to recipes
            // variable and focus on the newly created one
            formName = d3.select("#dropdownFormName").property('value');
            createFunc(formName);
        });
}

// renders the primary list under the management tab
function renderListItems() {
   /*
    <a href="#" class="list-group-item list-group-item-action py-3 lh-tight active ">
        <div class="col-10 mb-1 small">Some placeholder content in a paragraph below the heading and date.</div>
    </a>
    */

    query = d3.select("#input-search-query").property("value")
    items = store.listItems;
    if (query != "") {
        items = store.listItems.filter(function(d) {
            return d.name.toLowerCase().includes(query.toLowerCase());
        });
    }

    list = d3.select("#list-items");
    list.selectAll("a").remove();

    list.selectAll("a")
        .data(items)
        .enter()
        .append("a")
            .attr("class", "list-group-item list-group-item-action py-3 lh-tight")
            .attr("id", d => "list-" + d.name)
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
                "<div class=\"col-1\">"+(d.recipe_count > 0 ? d.recipe_count : "")+"</div>"+
                "</div>"
            });
}