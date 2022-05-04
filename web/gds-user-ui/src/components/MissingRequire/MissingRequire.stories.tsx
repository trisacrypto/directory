import { Meta, Story } from '@storybook/react';
import MissingRequire from '.';

type MissingRequireProps = {
  missingFields: Record<string, string | number | null>;
};

export default {
  title: 'components/MissingRequire',
  component: MissingRequire
} as Meta;

const Template: Story<MissingRequireProps> = (args) => <MissingRequire {...args} />;

export const Standard = Template.bind({});
Standard.args = {
  missingFields: {
    title: 'Please enter the title',
    phone_number: 'Please enter a valid phone number (e.g.: 0232 434 3242)',
    first_name: 'Please enter the first name',
    last_name: 'Please enter the last name'
  }
};
