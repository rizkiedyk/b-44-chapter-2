function submitData() {
  let name = document.getElementById("name").value;
  let email = document.getElementById("email").value;
  let phone = document.getElementById("phone").value;
  let subject = document.getElementById("subject").value;
  let message = document.getElementById("message").value;

  let emailReceiver = "akusiapa@gmail.com";

  if (name == "") {
    return alert("nama harus diisi");
  } else if (email == "") {
    return alert("email harus diisi");
  } else if (phone == "") {
    return alert("nomor telepon harus diisi");
  } else if (subject == "") {
    return alert("subject harus diisi");
  } else if (message == "") {
    return alert("message harus diisi");
  }

  let link = document.createElement("a");
  link.href = `mailto:${emailReceiver}?subject=${subject}&body=Hallo nama saya ${name}, ${message}, silahkan kontak nomer saya di ${phone}`;
  link.click();

  let student = {
    name,
    email,
    phone,
    subject,
    message,
  };

  console.log(student);
}
