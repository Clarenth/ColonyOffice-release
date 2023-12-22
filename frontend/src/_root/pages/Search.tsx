import React, { useState } from 'react'

// Components
import { Input } from '@/components/ui/input'
import SearchResults from '@/components/shared/SearchResults';
import GridDocsList from '@/components/shared/GridDocsList';


const Search = () => {
  const [searchValue, setSearchValue] = useState('')

  const docs = [];

  const shouldShowSearchResults = searchValue !== '';
  const shouldShowDocs = !shouldShowSearchResults && docs !== null;

  return (
    <React.Fragment>
      <div className='explore-container'>
        <div className='explore-inner_container'>
          <h2 className='h3-bold md:h2-bold w-full '>Search</h2>
          <div className='flex gap-1 px-4 w-full rounded-lg bg-dark-4'>
            <img 
              src='/assets/search-file-dark-svgrepo-com.svg'
              alt='search for files'
              width={24}
              height={24}
            />
            <Input 
              type='search'
              placeholder='Search'
              className='explore-search'
              value={searchValue}
              onChange={(event) => event.target.value}
            />
          </div>
        </div>

        <div className='flex-between w-full max-w-5xl mt-16'>
          <h3 className='body-bold md:h3-bold'>Newest Documents</h3>
          <div className='flex-center gap-3 bg-dark-3 rounded-xl px-4 py-2 cursor-pointer'>
            <p className='small-medium md:base-medium text-light-1'>all</p>
            <img 
              src='/assets/filter-dark-svgrepo-com.svg'
              height={24}
              width={24}
              alt='filter'
            />
          </div>
        </div>

        {/* <div className='flext flext-wrap gap-9 w-full max-w-5xl'>
          {
            shouldShowDocs ? 
            (
              <SearchResults />
            ) : shouldShowDocs ?
            (
              <p className='text-light-4 mt-10 text-center w-full'>Docs List</p>
            )
          : docs.pages.map((doc, key) => (
            <GridDocsList key={`doc-${key}`} doc={doc.}/>
          ))}
        </div> */}

      </div>
    </React.Fragment>
  )
}

export default Search