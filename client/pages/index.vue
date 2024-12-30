<template>
    <Card>
        <CardHeader>
            <h1 class="text-2xl font-bold">Welcome to your Nuxt.js app!</h1>
        </CardHeader>
        <CardContent class="p-4">
            <p>Click on the test button below to test!</p>
            <p>Test count: {{ testCount }}</p>
        </CardContent>
        <CardFooter>
            <Button class="w-full" @click="testCount++">Test</Button>
        </CardFooter>
    </Card>
</template>

<script setup>
const testCount = ref(0);
</script>