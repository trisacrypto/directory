import { Meta, Story } from '@storybook/react';
import LanguagesDropdown from 'components/LanguagesDropdown';

export default {
  title: 'components/LanguagesDropdown',
  component: LanguagesDropdown
} as Meta;

const Template: Story = (args) => <LanguagesDropdown {...args} />;

export const Standard = Template.bind({});
