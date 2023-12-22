import React, { useEffect } from 'react'

const AsyncRequests = () => {

  const login = () => {
    fetch("http://localhost:4000/auth/signin", 
    {
      method: 'POST',
      headers: 
      {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(
        {
          email: "",
          password: "",
        }
      ),
    })
    .then(response => response.json())
    .then(data => {
      sessionStorage.setItem("idToken", data.tokens.idToken)
      sessionStorage.setItem("refreshToken", data.tokens.refreshToken)
    })
    .catch(error => console.log(error))
  }
  
  const getAccount = async () => {
    const accountData = await fetch("http://localhost:4000/auth/account", 
    {
      method: 'GET',
      headers:
      {
        Authorization: `Bearer ${sessionStorage.getItem("idToken")}`
      }
    }
    )
    .then(response => response.json())
    .then((data) => {
      const account = data;
      //console.log("Hello from inside fetch: ", accountData)
      return account;
    })
    //console.log("Hello accountData: ", accountData)
    return accountData;
  }

  const accountPromise = async () => {
    const getAccount = await (await fetch("http://localhost:4000/auth/account", {
      method: 'GET',
      headers:
      {
        Authorization: `Bearer ${sessionStorage.getItem("idToken")}`
      }
    })).json()

    const data = Promise.all(
      [getAccount].map((obj) => {
        console.log(obj)
        const account = obj
        // const account = {
        //   firstName: obj.account.employee_identity_data.first_name,
        //   lastName: obj.account.employee_identity_data.last_name,
        //   email: obj.account.email,
        //   phoneNumber: obj.account.phone_number,
        //   jobTitle: obj.account.job_title,
        //   officeAddress: obj.account.office_address,
        //   employmentDate: obj.account.employment_date,
        // }
        return account
      })
    )
    return data
  }

  const getAccountData = async () => {
    try {
      const getAccount = await fetch("http://localhost:4000/auth/account", {
        method: 'GET',
        headers:
        {
          Authorization: `Bearer ${sessionStorage.getItem("idToken")}`
        }
      })
      const result = await getAccount.json();
      [result].map((data) => {
        console.log(data)
        const account = {
          firstName: data.employee_identity_data.first_name,
          middleName: data.employee_identity_data.middle_name,
          lastName: data.employee_identity_data.last_name,
          email: data.email,
          phoneNumber: data.phone_number,
          jobTitle: data.job_title,
          officeAddress: data.office_address,
          employmentDate: data.employment_date,
        }
        console.log("Hello from displayAccount: ", account)
        return account;
      })
    } catch (error) {
      console.log(error)
    }
  }

  const displayAccount = async() => {
    const result = await getAccount();
    console.log(result)
    return result;
    [result].map((data) => {
      //console.log(data)
      const account = {
        firstName: data.account.employee_identity_data.first_name,
        lastName: data.account.employee_identity_data.last_name,
        email: data.account.email,
        phoneNumber: data.account.phone_number,
        jobTitle: data.account.job_title,
        officeAddress: data.account.office_address,
        employmentDate: data.account.employment_date,
      }
      console.log("Hello from displayAccount: ", account)
      return account;
    })
  }

  useEffect(() => {
    //login();
    // setTimeout(() => { console.log("wait 2") }, 2000)
    //accountPromise();
    //console.log("getAccount: ", getAccount());
    //displayAccount();
    // const result = getAccount();
    // console.log(result)
    //displayAccount();
  })

  return (
    <div>Login</div>
  )
}

export default AsyncRequests;