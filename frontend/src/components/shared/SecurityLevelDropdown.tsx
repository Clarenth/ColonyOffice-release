// Libraries
import React from 'react'

// Components
import { DropdownMenu, DropdownMenuContent, DropdownMenuGroup, DropdownMenuItem, DropdownMenuTrigger } from '@radix-ui/react-dropdown-menu'
import { Button } from '../ui/button'

const SecurityLevelDropdown = () => {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Button variant="outline">Security Access Level</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className='w-2'>
        <DropdownMenuGroup>
          <DropdownMenuItem>
            <span>Official</span>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <span>Secret</span>
          </DropdownMenuItem>
          <DropdownMenuItem>
            <span>TopSecret</span>
          </DropdownMenuItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

export default SecurityLevelDropdown