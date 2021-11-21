import './App.scss';
import React, {useState} from 'react'
import {ToastContainer, toast} from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css';

const initialValues={
  name:null,
  email:null,
  phone:null,
  message:null
}
function App() {
const [state,setState]= useState(initialValues)

const {name,email,phone,message} = state
 

const handleSubmit = (e) => {
  e.preventDefault();
  if (name === null) {
    toast.error("Please provide your Name");
  }else if(email === null){
    toast.error("Please provide your Email");
  } else if(phone === null){
    toast.error("Please provide your Phone Number");
  } 
  else {
    fetch('http://localhost:8080/submit-form', {
      method: 'POST',
      redirect: 'follow',
      body: JSON.stringify({
        name,
        email,
        phone,
        message
      }),
      headers: {
        'Content-Type': 'application/json',
      }
    }).then(resp => {
      if (resp.status >= 200 && resp.status < 400) {
        console.log(resp);
        if (resp.redirected) {
          window.location.href = resp.url
        }
        return resp;
      } else {
        console.log("Form submission failed ", resp.status, resp.statusText)
      }
    }).catch(e => console.log("POST form request failed", e))
   
    toast.success("Form Submitted Successfully");
  }
};

  const handleInputChange=(e)=>{
    let {name,value} = e.target
    setState({...state, [name]:value})
  }
  console.log(name,'name')
  console.log(email,'email')
  return (
   <section className="contact-section" >
    <div className='container' >
    <ToastContainer position='top-center'/>
     <div className='row justify-content-center' >
    <div className='col-md-10'>
    <div className='wrapper' >
    <div className='row no-gutters' >
    <div className='col-md-6' >
    <div className='contact-wrap w-100 p-lg-5 p-4'>
    <h3 className='mb-4'>Send us Your Info</h3>
    <form id='contact-form' className='contactForm' onSubmit={handleSubmit} >
    <div className='row'>
    <div className='col-md-12'>
    <div className='form-group' >
    <input 
      type='text'
      className='form-control'
      name='name'
      placeholder='Name'
      onChange={handleInputChange}
      value={name}
      
    />

    </div>

    </div>
    <div className='col-md-12'>
    <div className='form-group' >
    <input 
      type='email'
      className='form-control'
      name='email'
      placeholder='Email'
      onChange={handleInputChange}
      value={email}
      

      
    />

    </div>

    </div>
    <div className='col-md-12'>
    <div className='form-group' >
    <input 
      type='text'
      className='form-control'
      name='phone'
      placeholder='Phone Number'
      onChange={handleInputChange}
      value={phone}
     
    />

    </div>

    </div>
    <div className='col-md-12'>
    <div className='form-group' >
    <textarea 
      type='text'
      className='form-control'
      name='message'
      placeholder='Message'
      cols='30'
      rows='6'
      onChange={handleInputChange}
      value={message}
    />

    </div>

    </div>
    <div className='col-md-12' >
    <div className='form-group' >
    <input 
      type='submit'
      value='Send Info'
      className='btn btn-primary'
    /> 
    </div>

    </div>
    </div>
    
    </form>
    </div>

    </div>

<div className='col-md-6 d-flex align-items-stretch'>
<div className='info-wrap w-100 p-lg-5 p-4' >
<h3>Contact us</h3>
<p className='mb-4'  >Leave your Info, Leave a suggestion, Fumar Mota!</p>
<div className='dbox w-100 d-flex align-items-start' >
<div className='icon d-flex align-items-center justify-content-center' ></div>
<span className='fa fa-map-marker' ></span>
<div className='text pl-3' >
  <p>
    <span>_  Address:</span> Far-Far AF, Columbia Suite 420 gang-gang
  </p>
</div>

</div>
<div className='dbox w-100 d-flex align-items-start' >
<div className='icon d-flex align-items-center justify-content-center' ></div>
<span className='fa fa-phone' ></span>
<div className='text pl-3' >
  
    <span>_  Phone:</span>
    <a href='tel://99999999'>+33 42069 333</a>
   
</div>

</div>
<div className='dbox w-100 d-flex align-items-start' >
<div className='icon d-flex align-items-center justify-content-center' ></div>
<span className='fa fa-paper-plane' ></span>
<div className='text pl-3' >
  
    <span>_  Email:</span> 
    <a href='mailto:walkerbishop333@gmail.com'> ColumbiaChronicIreland.org</a>
  
</div>

</div>
<div className='dbox w-100 d-flex align-items-start' >
<div className='icon d-flex align-items-center justify-content-center' ></div>
<span className='fa fa-globe' ></span>
<div className='text pl-3' >
  
    <span>_ Website:</span> <a href='#'>Your_Main_Website.com</a>
  
</div>

</div>
</div>
</div>

    </div>
    </div>

    </div>
     </div>
    </div>
   </section>
  );
}

export default App;
