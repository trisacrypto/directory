import React, { FormEvent } from 'react';
import { Story } from '@storybook/react';
import SearchDirectory from './SearchDirectory';

interface SearchDirectoryProps {
  handleSubmit: (e: FormEvent, query: string) => void;
  isLoading: boolean;
  result: any;
  error: string;
  query: string;
  options: any;
  onResetData: () => void;
}

export default {
  title: 'Components/SearchDirectory',
  component: SearchDirectory
};

export const standard: Story<SearchDirectoryProps> = ({ ...props }) => (
  <SearchDirectory {...props} />
);

standard.bind({});
