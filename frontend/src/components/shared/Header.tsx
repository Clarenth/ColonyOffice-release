// Libraries
import React, { useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'

// Mutators
import { useLogoutAccountMutation } from '@/lib/react-query/queriesAndMutations'

// Components
import { Button } from '../ui/button'
import { useAccountContext } from '@/context/AuthContext'

const Header = () => {
  // Hooks
  const navigate = useNavigate();
  const { mutateAsync: logout, isSuccess } = useLogoutAccountMutation();

  // Context
  const { account } = useAccountContext();

  useEffect(() => {
    if(isSuccess) {
      navigate(0)
    }
  }, [isSuccess])

  return (
    <section className='header'>
      <div className='flex-between py-4 px-5'>
        <Link to="/" className='flex gap-3 items-center'>
          logo goes here
          {/* <img
            src='logo goes here'
            width={130}
            height={325}
          /> */}
        </Link>
        <div className='flex gap-4'>
          <Button onClick={() => logout} variant="ghost" className='shad-button_ghost'>
            logout
            <img />
          </Button>
          <Link to={`/account/${account}`} className='flex-center gap-3'>
            <img 
              src={'/assets/images/profile-placeholder.svg'}
            />
          </Link>
        </div>
      </div>
    </section>
  )
}

export default Header