// Libraries
import { Link, useNavigate } from "react-router-dom"
import { useForm } from "react-hook-form"
import * as zod from "zod"
import { zodResolver } from "@hookform/resolvers/zod"

// Validation
import { SignupValidation } from "@/lib/validation"

// React-Query Mutations
import { useCreateAccountMutation, useLoginAccountMutation } from "@/lib/react-query/queriesAndMutations";

// Context
import { useAccountContext } from "@/context/AuthContext"

// Components
import { Button } from "@/components/ui/button"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage, } from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { useToast } from "@/components/ui/use-toast"
// Shared
import Loader from "@/components/shared/Loader"

const Signup = () => {
  // Hooks
  const navigate = useNavigate();
  const { toast } = useToast();
  const { checkAuthAccount, isLoading } = useAccountContext();

  const { mutateAsync: createAccount, isPending: isCreatingAccount } = useCreateAccountMutation();
  const { mutateAsync: loginAccount, isPending: isLoggingIn } = useLoginAccountMutation();
  
  // 1. Define your form.
  const form = useForm<zod.infer<typeof SignupValidation>>({
    resolver: zodResolver(SignupValidation),
    defaultValues: {
      email: "",
      password: "",
      phone_number: "",
      job_title: "",
      office_address: "",
      employee_identity_data:
      {
        first_name: "",
        middle_name: "",
        last_name:"",
        sex: "",
        gender: "",
        age:"",
        height:"",
        home_address:"",
        birthdate:"",
        birthplace:"",
      },
      security_access_level: "",
    },
  })
 
  // 2. Define a submit handler.
  async function onSubmit(values: zod.infer<typeof SignupValidation>) {
    const newAccount = await createAccount(values);
    console.log(newAccount);
    if(!newAccount){
      return toast({
        title: "Signup failed! Please try again.",
      })
    }

    // Create session is likely not going to be done. Should navigate to login
    const session = await loginAccount({
      email: values.email,
      password: values.password,
    });
    if(!session) {
      return toast({ title: 'Login failed. Please try again.' })
    }

    const isLoggedIn = await checkAuthAccount();

    if(isLoggedIn) {
      form.reset()
      navigate("/")
    } else {
      return toast({ title: "Login failed. Please try again." })
    }
  }
  
  return (
    <Form {...form}>
      <div className="sm:w-420 flex-center flex-col"> {/*</div><div className="sm:w-420 flex-center flex-col">*/}
        <img src="/assets/images/logo.svg" alt="logo" />
        <h2 className="h3-bold md:h2-bold pt-5 sm:pt-12">Create Account</h2>
        <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-2 w-full mt-4">
          <div className="flex flex-row gap-3 w-full mt-4">
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Email</FormLabel>
                <FormControl>
                  <Input 
                    type="email" 
                    className="shad-input"
                    placeholder="Email"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Password</FormLabel>
                <FormControl>
                  <Input 
                    type="password" 
                    className="shad-input"
                    placeholder="Password"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          </div>
          <div className="flex flex-row gap-3 w-full mt-4">
          <FormField
            control={form.control}
            name="phone_number"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Phone Number</FormLabel>
                <FormControl>
                  <Input 
                    type="tel" 
                    className="shad-input"
                    placeholder="Phone Number"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="job_title"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Job Title</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Job Title"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="office_address"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Office Address</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Office Address"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          </div>
          <div className="flex flex-row gap-3 w-full mt-4">
          <FormField
            control={form.control}
            name="employee_identity_data.first_name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>First Name</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="First Name"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="employee_identity_data.middle_name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Middle Name</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Middle Name"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="employee_identity_data.last_name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Last Name</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Last Name"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          </div>
          <div className="flex flex-row gap-3 w-full mt-4">
          <FormField
            control={form.control}
            name="employee_identity_data.sex"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Sex</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Sex"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="employee_identity_data.gender"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Gender</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Gender"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="employee_identity_data.age"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Age</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Age"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="employee_identity_data.height"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Height</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Height"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          </div>
          <div className="flex flex-row gap-3 w-full mt-4">
          <FormField
            control={form.control}
            name="employee_identity_data.home_address"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Home Address</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Home Address"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="employee_identity_data.birthdate"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Birthdate</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Birthdate"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="employee_identity_data.birthplace"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Birthplace</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Birthplace"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          </div>
          <div className="flex-row gap-3 w-full mt-4">
          <FormField
            control={form.control}
            name="security_access_level"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Security Access Level</FormLabel>
                <FormControl>
                  <Input 
                    type="text" 
                    className="shad-input"
                    placeholder="Security Access Level"
                    {...field} 
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          </div>
          <Button type="submit" className="mt-4 shad-button_primary">
            { isCreatingAccount ? (
              <div className="flex-center gap-2">
                <Loader />
                Loading...
              </div>
            ): "Signup"}
          </Button>
          <p className="text-small-regular text-light-2 text-center mt-2">
              Already have an account?
              <Link to="/login" className="text-primary-500 text-small-semibold ml-1">Login</Link>
          </p>
        </form>
      </div>
    </Form>
  )
}

export default Signup;