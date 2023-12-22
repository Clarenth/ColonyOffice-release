// Libraries
import React, { useEffect } from 'react'
import { Link, NavLink, useLocation, useNavigate } from 'react-router-dom'

// Mutators
import { useLogoutAccountMutation } from '@/lib/react-query/queriesAndMutations'

// Components
import { Button } from '../ui/button'
import { useAccountContext } from '@/context/AuthContext'
import { sidebarLinks } from '@/constants/index'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '../ui/tooltip'

const SidebarLeft = () => {
  // Hooks
  const { pathname } = useLocation();
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
    <nav className='leftsidebar'>
      <div className='flex flex-col gap-11 '>
        <Link to={`/profile/${account.email}`} className='flex gap-3 items-center'>
          <img
            src='/assets/collaboration-dark-svgrepo-com.svg'
            width={24}
            height={24}
            className='h-14 w-14 rounded-full'
          />
          <div className='flex flex-col'>
            <p className='body-bold'>
              {account.firstName + " " + account.lastName}
            </p>
            <p className='samll-regular text-light-3'>
              {account.email}
            </p>
          </div>
        </Link>
        <ul>
          { sidebarLinks.map((link) => {
            const isActive = pathname === link.route
            return (
              <li key={link.label} className={`leftsidebar-link group ${isActive && 'bg-blue-800'}`}>
                {/* <TooltipProvider> */}
                  {/* <Tooltip> */}
                    {/* <TooltipTrigger asChild> */}
                      <NavLink
                        to={link.route}
                        className="flex gap-4 items-center p-4"
                      >
                        <img 
                          src={link.icon} 
                          alt={link.label}
                          width={36}
                          height={36}
                          className={`group-hover:invert-white ${isActive && 'invert-white'}`}
                        />
                        {link.label}
                      </NavLink>
                    {/* </TooltipTrigger> */}
                    {/* <TooltipContent > */}
                      {/* <p>{link.label}</p> */}
                    {/* </TooltipContent> */}
                  {/* </Tooltip> */}
                {/* </TooltipProvider> */}
              </li>
            )
          })}
        </ul>
      </div>

      <Button onClick={() => logout} variant="ghost" className='shad-button_ghost'>
        <img 
          src='/assets/sign-out-dark-svgrepo-com.svg'
          alt='logout'
          width={36}
          height={36}
        />
        <p className='small-medium lg:base-medium'>Logout</p>
      </Button>
    </nav>
  )
}

export default SidebarLeft