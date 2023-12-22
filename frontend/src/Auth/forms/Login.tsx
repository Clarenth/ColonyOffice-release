// Libraries
import { Link, useNavigate } from "react-router-dom"
import { useForm } from "react-hook-form"
import * as zod from "zod"
import { zodResolver } from "@hookform/resolvers/zod"

// Internal Lib
import { LoginValidation } from "@/lib/validation"

// React-Query Mutations
import { useLoginAccountMutation } from "@/lib/react-query/queriesAndMutations";

// API
import { postLoginAccount } from "@/lib/colony-office/api"

// Context
import { useAccountContext } from "@/context/AuthContext"

// Components
import { Button } from "@/components/ui/button"
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage, } from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { useToast } from "@/components/ui/use-toast"
// Shared
import Loader from "@/components/shared/Loader"

const Login = () => {
  // Hooks
  const navigate = useNavigate();
  const { toast } = useToast();
  const { checkAuthAccount, isLoading } = useAccountContext();

  const { mutateAsync: loginAccount } = useLoginAccountMutation();
  
  // 1. Define your form.
  const form = useForm<zod.infer<typeof LoginValidation>>({
    resolver: zodResolver(LoginValidation),
    defaultValues: {
      email: "",
      password: "",
    },
  })
 
  // 2. Define a submit handler.
  async function onSubmit(values: zod.infer<typeof LoginValidation>) {
    const newAccount = await postLoginAccount(values);
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
        <h2 className="h3-bold md:h2-bold pt-5 sm:pt-12">Login</h2>
        <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-2 w-full mt-4">
          <div className="flex flex-col gap-3 w-full mt-4">
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
          <Button type="submit" className="mt-4 shad-button_primary">
            { isLoading ? (
              <div className="flex-center gap-2">
                <Loader />
                Loading...
              </div>
            ): "Login"}
          </Button>
          <p className="text-small-regular text-light-2 text-center mt-2">
              No account?
              <Link to="/signup" className="text-primary-500 text-small-semibold ml-1">Signup!  </Link>
          </p>
        </form>
      </div>
    </Form>
  )
}

export default Login;