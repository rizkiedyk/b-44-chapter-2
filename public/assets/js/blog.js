let blogs = [];

function getData(event) {
  event.preventDefault();

  let name = document.getElementById("project-name").value;
  let startDate = new Date(document.getElementById("start-date").value);
  let endDate = new Date(document.getElementById("end-date").value);
  let description = document.getElementById("description-box").value;
  let image = document.getElementById("input-image").files;

  image = URL.createObjectURL(image[0]);

  // technologies checked

  let nodeJs = document.getElementById("nodejs");
  let reactJs = document.getElementById("reactjs");
  let nextJs = document.getElementById("nextjs");
  let typeScript = document.getElementById("typescript");

  let iconNodeJs = ``;
  let iconReactJs = ``;
  let iconNextJs = ``;
  let iconTypeScript = ``;

  if (nodeJs.checked === false) {
    iconNodeJs = `style="display: none;"`;
  }
  if (reactJs.checked === false) {
    iconReactJs = `style="display: none;"`;
  }
  if (nextJs.checked === false) {
    iconNextJs = `style="display: none;"`;
  }
  if (typeScript.checked === false) {
    iconTypeScript = `style="display: none;"`;
  }

  let addBlog = {
    name,
    startDate,
    endDate,
    description,
    image,
    iconNodeJs,
    iconReactJs,
    iconNextJs,
    iconTypeScript,
  };

  console.log(addBlog);
  blogs.push(addBlog);
  showData();
}

function showData() {
  document.getElementById("post-blog").innerHTML = ``;
  for (let i = 0; i < blogs.length; i++) {
    document.getElementById("post-blog").innerHTML += `<div class="card-blog">
    <div class="content-blog">
      <div class="thumbnail">
        <img src="${blogs[i].image}" alt="" />
      </div>

      <div class="title-blog">
        <h3><a href="blog.html">${blogs[i].name}</a></h3>
        <p>durasi : ${getDistanceTime(blogs[i].endDate, blogs[i].startDate)}</p>
      </div>

      <div class="description-output">
        ${blogs[i].description}
      </div>

      <div class="technologies-output">
        <div class="nodejs">
          <img src="assets/img/nodejs.png" alt="" ${blogs[i].iconNodeJs} />
        </div>
        <div class="nextjs">
          <img src="assets/img/nextjs.png" alt="" ${blogs[i].iconNextJs}/>
        </div>
        <div class="reactjs">
          <img src="assets/img/reactjs.png" alt="" ${blogs[i].iconReactJs}/>
        </div>
        <div class="typescript">
          <img src="assets/img/typescript.png" alt="" ${
            blogs[i].iconTypeScript
          }/>
        </div>
      </div>

      <div class="opsi-blog">
        <div class="edit-blog">
          <a href="">edit</a>
        </div>

        <div class="delete-blog">
          <a href="">delete</a>
        </div>
      </div>
    </div>
  </div>`;
  }
}

// // 1 tahun = 12
// // 1 bulan = 30
// // 1 hari = 24
// // 1 jam = 60
// // 1 menit = 60
// // 1 detik = 1000

function getDistanceTime(endDate, startDate) {
  let distance = endDate - startDate;

  let distanceMonth = Math.floor(distance / 1000 / 60 / 60 / 24 / 30);
  let distanceDay = Math.floor(distance / 1000 / 60 / 60 / 24);

  if (distanceMonth > 0) {
    if (distanceDay % 30 >= 1) {
      return `${distanceMonth} Bulan ${distanceDay % 30} Hari`;
    }
    return `${distanceMonth} month`;
  } else if (distanceDay > 0) {
    return `${distanceDay} Hari`;
  } else {
    return `0 Hari`;
  }
}
