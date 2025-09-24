<script lang="ts">
	import * as Form from '$lib/components/ui/form';
	import { m } from '$lib/paraglide/messages';
    
    let {
        fields = $bindable(),
        errors,
        disabled
    } : {
		fields: any;
		errors?: any;
        disabled?: boolean;
    } = $props()
</script>

<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
    <Form.Checkbox
        title={m.courseitem_group_pick()}
        name="choice"
        id="choice"
        bind:value={fields.choice}
        error={errors.choice}
        {disabled}
    ></Form.Checkbox>
    <Form.TextInput
        title={m.courseitem_group_pick_min()}
        name="chooseMin"
        id="chooseMin"
        type="number"
        bind:value={fields.chooseMin}
        error={errors.chooseMin ?? ''}
        disabled={!fields.choice || disabled}
    ></Form.TextInput>
    <Form.TextInput
        title={m.courseitem_group_pick_max()}
        name="chooseMax"
        id="chooseMax"
        type="number"
        bind:value={fields.chooseMax}
        error={errors.chooseMax ?? ''}
        disabled={!fields.choice || disabled}
    ></Form.TextInput>
</div>