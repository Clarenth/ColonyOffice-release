import React from 'react'
import { Link, useLocation } from 'react-router-dom'

//
import { footerLinks } from '@/constants'

const Footer = () => {
  const { pathname } = useLocation();

  return (
    <section className='bottom-bar'>
      { footerLinks.map((link) => {
        const isActive = pathname === link.route
        return (
            <Link
              to={link.route}
              key={link.label} className={`${isActive && 'bg-blue-800 rounded-[10px] flex-center flex-col gap-1 p-2 transition'}`}
            >
              <img 
                src={link.icon} 
                alt={link.label}
                width={20}
                height={20}
                className={`${isActive && 'invert-white'}`}
              />
              {/* <p className='tiny-medium text-light-2'>{link.label}</p> */}
            </Link>
        )
      })}
    </section>
  )
}

export default Footer