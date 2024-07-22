import { Meta, Story } from '@storybook/react';
import { withRHF } from 'hoc/withRHF';
import OtherJurisdictions from '.';

type OtherJurisdictionsProps = {
  name: string;
};

export default {
  title: 'components/OtherJurisdictions',
  component: OtherJurisdictions,
  decorators: [withRHF(false)]
} as Meta;

const Template: Story<OtherJurisdictionsProps> = (args) => <OtherJurisdictions {...args} />;

export const Standard = Template.bind({});
Standard.args = {
  name: 'trixo.other_jurisdictions'
};
