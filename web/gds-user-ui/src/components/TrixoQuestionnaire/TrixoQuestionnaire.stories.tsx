import { Meta, Story } from '@storybook/react';
import TrixoQuestionnaire from '.';

export default {
  title: 'components/TrixoQuestionnaire',
  component: TrixoQuestionnaire
} as Meta;

const Template: Story = (args) => <TrixoQuestionnaire {...args} />;

export const Standard = Template.bind({});
Standard.args = {};
