import React, { FormEvent, useState } from 'react';

import Head from 'components/Head/LandingHead';
import JoinUsSection from 'components/Section/JoinUs';
import SearchDirectory from 'components/Section/SearchDirectory';
import AboutTrisaSection from 'components/Section/AboutUs';

import LandingLayout from 'layouts/LandingLayout';
import useFetchLookupAutocomplete from './useFetchLookupAutocomplete';
import useFetchLookup from './useFetchLookup';
import NetworkActivity from 'components/Section/NetworkActivity/NetworkActivity';
import { Hide } from '@chakra-ui/react';
import LandingBanner from 'components/Banner/LandingBanner';
const HomePage: React.FC = () => {
  const { vasps } = useFetchLookupAutocomplete();
  const { handleSearch, searchString, data, isLoading, error, resetData } = useFetchLookup();
  const [search, setSearch] = useState('');
  const handleSearchSubmit = (evt: FormEvent, searchQuery: string) => {
    evt.preventDefault();
    handleSearch(searchQuery);
    setSearch(searchString);
  };

  return (
    <LandingLayout>
      <LandingBanner />
      <Head hasBtn isHomePage />
      <AboutTrisaSection />
      <JoinUsSection />
      <Hide>
        <NetworkActivity />
      </Hide>

      <SearchDirectory
        handleSubmit={handleSearchSubmit}
        isLoading={isLoading}
        result={data}
        error={error}
        handleClose={() => resetData()}
        onResetData={() => resetData()}
        query={search}
        options={vasps}
      />
    </LandingLayout>
  );
};

export default HomePage;
