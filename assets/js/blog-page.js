let blogs = [];

function getData(event) {
    event.preventDefault()

    let projectName = document.getElementById('project-name').value;
    let startDate = new Date(document.getElementById('start-date').value);
    let endDate = new Date(document.getElementById('end-date').value);
    let description = document.getElementById('description-box').value;
    // let technologies = document.getElementById('').value;
    let image = document.getElementById('input-image').files;

    image = URL.createObjectURL(image[0]);

    let addBlog = {
        projectName,
        startDate,
        endDate,
        description,
        image,
        postedAt : new Date ()
    }

    blogs.push(addBlog);

    console.log(blogs);
    showData();
}

function showData() {
    document.getElementById("post-blog").innerHTML = ``
    for( i = 0; i < blogs.length; i++) {
        document.getElementById('post-blog').innerHTML += `<div class="card-blog">

        <div class="content-blog">

            <div class="thumbnail">
                <img src="${blogs[i].image}" alt="">
            </div>
            
            <div class="title-blog">
                <h3><a href="blog.html">${blogs[i].projectName}</a></h3>
                <p>durasi : ${getDistanceTime(blogs[i].endDate, blogs[i].startDate)}</p>
            </div>
        
            <div class="description-output">
            ${blogs[i].description}    
            </div>

            <div class="technologies-output">
                <div class="nodejs">
                    <img src="assets/img/nodejs.png" alt="">
                </div>
                <div class="nextjs">
                    <img src="assets/img/nextjs.png" alt="">
                </div>
                <div class="reactjs">
                    <img src="assets/img/reactjs.png" alt="">
                </div>
                <div class="typescript">
                    <img src="assets/img/typescript.png" alt="">
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

    </div>`
    }
}

// // 1 tahun = 12
// // 1 bulan = 30
// // 1 hari = 24
// // 1 jam = 60
// // 1 menit = 60
// // 1 detik = 1000

function getDistanceTime(endDate, startDate){
    
    let distance = endDate - startDate

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
