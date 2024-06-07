import { A, action, createAsync, redirect, useAction, useSubmission } from "@solidjs/router";
import { For, Show, Suspense, createEffect, createSignal } from "solid-js";
import Counter from "~/components/Counter";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { TextField, TextFieldInput, TextFieldLabel, TextFieldTextArea } from "~/components/ui/text-field";
import createArticle from "~/lib/rpc/createArticle";
import getArticles from "~/lib/rpc/getArticles";
import { Title } from "@solidjs/meta"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "~/components/ui/collapsible";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "~/components/ui/tabs";
import searchSemantics from "~/lib/rpc/searchSemantics";
import { Alert, AlertDescription, AlertTitle } from "~/components/ui/alert";
import { Article, ArticleWithScore } from "~/lib/rpc";


export default function Home() {
  return (
    <main class="max-w-3xl m-auto pt-8">
      <Title>Verctor Search</Title>
      <Tabs defaultValue="search" class="w-full grid place-items-center">
        <TabsList class="mx-auto">
          <TabsTrigger value="search">search</TabsTrigger>
          <TabsTrigger value="create">create</TabsTrigger>
          <TabsTrigger value="list">list</TabsTrigger>
        </TabsList>
        <div class="w-full">
          <TabsContent value="search">
            <SemanticsSearch />
          </TabsContent>
          <TabsContent value="create">
            <ArticleCreator />
          </TabsContent>
          <TabsContent value="list">
            <Articles />
          </TabsContent>
        </div>

      </Tabs>
    </main>

  );
}


function Articles() {
  const articles = createAsync(getArticles)
  return (
    <For each={articles()} fallback={<div>Loading...</div>}>
      {article => <ArticleCard article={article} />}
    </For>
  )
}

const createAction = action(createArticle)

function ArticleCreator() {
  const [title, setTitle] = createSignal('')
  const [content, setContent] = createSignal('')
  return (
    <Card class="m-4">
      <CardHeader>Create Article</CardHeader>
      <CardContent class="grid grid-cols-1 gap-4">
        <div>
          <TextField>
            <TextFieldLabel for="title">title</TextFieldLabel>
            <TextFieldInput type="text" id="title" value={title()} onInput={e => setTitle(e.currentTarget.value)} />
          </TextField>
        </div>
        <div>
          <TextField>
            <TextFieldLabel for="content">content</TextFieldLabel>
            <TextFieldTextArea rows={10} id="content" value={content()} onInput={e => setContent(e.currentTarget.value)} />
          </TextField>
        </div>
        <form action={createAction.with(title(), content())} method="post">
          <Button type="submit">Create</Button>
        </form>
      </CardContent>
    </Card>
  )
}

const searchAction = action(searchSemantics)
function SemanticsSearch() {
  const [query, setQuery] = createSignal('')
  const search = useAction(searchAction)
  const searching = useSubmission(searchAction)

  const handleSearch = async () => {
    search(query())
  }

  return (
    <div>
      <TextField>
        <TextFieldLabel for="query">query</TextFieldLabel>
        <TextFieldInput type="text" id="query" value={query()} onInput={e => setQuery(e.currentTarget.value)} />
      </TextField>
      <Button onClick={handleSearch}>Search</Button>
      <Show when={searching.pending}>
        <div>Loading...</div>
      </Show>
      <Show when={searching.error}>
        <Alert variant="destructive" class="my-3">
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>
            {searching.error.message}
          </AlertDescription>
        </Alert>
      </Show>
      <For each={searching.result}>
        {({ Article: article, Score: score }) => <ArticleCard article={article} score={score} />}
      </For>
    </div>
  )
}

type ArticleCardProps = {
  article: Article;
  score?: number;
}
function ArticleCard(props: ArticleCardProps) {
  if (props.score) {
    return (
      <Card class="m-4">
        <CardHeader>
          <CardTitle class="flex justify-between">
            <span>{props.article.title}</span>
            <span class="ml-auto mr-0 text-sm font-bold">score: {props.score}</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {props.article.content}
        </CardContent>
      </Card>
    )
  }
  return (
    <Card class="m-4">
      <CardHeader>
        <CardTitle>{props.article.title}</CardTitle>
      </CardHeader>
      <CardContent>
        {props.article.content}
      </CardContent>
    </Card>
  )
}