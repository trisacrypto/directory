import { Meta, Story } from '@storybook/react';
import { withRHF } from 'hoc/withRHF';
import OtherJuridictions from '.';

type OtherJuridictionsProps = {
  name: string;
};

export default {
  title: 'components/OtherJuridictions',
  component: OtherJuridictions,
  decorators: [withRHF(false)]
} as Meta;

const Template: Story<OtherJuridictionsProps> = (args) => <OtherJuridictions {...args} />;

export const Standard = Template.bind({});
Standard.args = {
  name: 'trixo.other_jurisdictions'
};
