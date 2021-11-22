<template>
  <el-container class="home-container">
    <div class="home-header"></div>
    <el-main class="home-main">
      <!-- 搜索框 -->
      <el-row>
        <el-col :span="1"></el-col>
        <el-col :xs="22" :sm="12" :md="9" :lg="9" :xl="9">
          <el-input
            autofocus
            v-model="query"
            autocomplete="off"
            placeholder=""
            v-on:keydown="enterSearch"
            ref="searchInput"
          >
            <template #append>
              <el-button
                @click="clickSearchButton(0, 0, 0, 0)"
                icon="el-icon-search"
              ></el-button></template
          ></el-input>
        </el-col>
      </el-row>

      <!-- 未输入搜索结果时的展示页 -->
      <el-row v-if="!searched" style="height: 40px">
        <el-col :span="1"></el-col>
        <el-col :span="11">
          <h3>搜索表达式样例：</h3>
        </el-col>
      </el-row>
      <el-row v-if="!searched">
        <el-col :span="1"></el-col>
        <el-col :span="11">
          <div class="col-md-8">
            <dl class="dl-horizontal">
              <dt>cpu</dt>
              <dd>搜包含 "cpu" 的文档</dd>
              <dt>"CPU cache"</dt>
              <dd>搜包含短语 "CPU cache" 的文档</dd>
              <dt>cpu.*name</dt>
              <dd>搜匹配正则表达式 "cpu.*name" 的文档</dd>
              <dt>"cpu\d{3} name"</dt>
              <dd>搜匹配正则表达式 "cpu\d{3} name" 的文档</dd>
              <dt>cpu cache</dt>
              <dd>搜包含 "cpu" 并且包含 "cache" 的文档</dd>
              <dd>同 cpu AND cache</dd>
              <dd>同 cpu and cache</dd>
              <dd>同 (cpu and cache)</dd>
              <dt>cpu or cache</dt>
              <dd>搜包含 "cpu" 或者 "cache" 的文档</dd>
              <dd>同 cpu OR cache</dd>
              <dd>同 (cpu or cache)</dd>
              <dt>cpu -cache</dt>
              <dd>搜包含 "cpu" 并且不包含 "cache" 的文档</dd>
              <dt>cpu -(cache or miss)</dt>
              <dd>搜包含 "cpu" 并且不包含 "cache" 或 "miss" 的文档</dd>
              <dt>cpu -"cache miss"</dt>
              <dd>搜索包含 "cpu" 但不含 "cache miss" 短语的文档</dd>
              <dt>cpu cache or hit miss</dt>
              <dd>等价于搜 (cpu AND cache) OR (hit AND miss)</dd>
              <dd>
                - 优先级高于 AND 高于
                OR，可以用这三个操作符加括号组合任意深度的表达式
              </dd>
            </dl>
          </div>
        </el-col>
        <el-col :span="11">
          <div class="col-md-8">
            <dl class="dl-horizontal">
              <dt>CPU cache.*name case:yes</dt>
              <dd>
                搜包含 "CPU" 并且匹配 "cache.*name"
                的文档，两者都大小写敏感。默认不区分大小写
              </dd>
              <dt>cpu file:admin</dt>
              <dd>搜包含 "cpu" 并且文件名中包含 "admin" 的文档</dd>
              <dt>cpu (file:api.*doc or file:admin)</dt>
              <dd>
                搜包含 "cpu" 并且文件名匹配正则表达式 "api.*doc" 或者包含
                "admin" 的文档
              </dd>
              <dt>cpu lang:java</dt>
              <dd>搜包含 "cpu" 的 java 代码</dd>
              <dt>cpu -lang:java</dt>
              <dd>搜包含 "cpu" 的非 java 代码</dd>
              <dt>cpu -(lang:java or lang:python)</dt>
              <dd>搜包含 "cpu" 的 java/python 之外的代码</dd>
              <dt>file:\.cpp$</dt>
              <dd>搜索文件名以 ".cpp" 结尾的文档</dd>
              <dt>cpu -file:admin.*java -file:web</dt>
              <dd>
                搜包含 "cpu" 但文件名不匹配正则表达式 "admin.*java" 也不包含
                "web" 的文档
              </dd>
              <dt>sym:data</dt>
              <dd>搜索包含符号 "data" 的文档，sym:不可以作用在正则表达式上</dd>
              <dt>cpu.*name repo:web.*service</dt>
              <dd>
                搜索匹配正则表达式 "cpu.*name"，同时所在仓库名匹配正则表达式
                "web.*service" 的文档
              </dd>
              <dt>repo:web.*service</dt>
              <dd>搜索仓库名匹配正则表达式 "web.*service" 的代码仓库</dd>
              <dt>repo:web.*service -repo:admin</dt>
              <dd>
                搜索仓库名匹配正则表达式 "web.*service" 但名称中不包含 "admin"
                的代码仓库
              </dd>
            </dl>
          </div>
        </el-col>
      </el-row>

      <!-- 搜索结果提示 -->
      <el-row v-if="searched">
        <el-col :span="1"></el-col>
        <el-col :span="22">
          <div class="tips">
            {{ message }}
          </div>
        </el-col>
      </el-row>

      <!-- 搜索结果 -->
      <el-row class="search-results">
        <el-col :span="1"></el-col>

        <!-- 左侧导航 -->
        <el-col :span="6" v-if="repos != null && repos.length > 0">
          <!-- 仓库列表 -->
          <div class="repos-list-title">仓库列表</div>
          <div class="repos-list">
            <div
              class="repo-list-item"
              v-for="repo in repos"
              :key="repo.repoID"
            >
              <a
                href="#"
                @click="
                  hasRepo(repo.remoteURL)
                    ? removeRepo(
                        repo.remoteURL == '' ? repo.localPath : repo.remoteURL
                      )
                    : addRepo(
                        repo.remoteURL == '' ? repo.localPath : repo.remoteURL
                      )
                "
              >
                <font-awesome-icon
                  v-if="!hasRepo(repo.remoteURL)"
                  icon="plus"
                ></font-awesome-icon>
                <font-awesome-icon
                  v-if="hasRepo(repo.remoteURL)"
                  icon="minus"
                ></font-awesome-icon>
                {{ simplifyRepo(repo.remoteURL) }}</a
              >
            </div>
          </div>

          <!-- 语言列表 -->
          <div
            v-if="
              searchType == 'documents' && langs != null && langs.length > 0
            "
          >
            <div class="files-list-title">语言列表</div>
            <div class="files-list">
              <div
                class="file-list-item"
                v-for="lang in langs"
                :key="lang.languageID"
              >
                <div v-if="lang.name != ''">
                  <a
                    href="#"
                    @click="
                      hasLang(lang.name)
                        ? removeLang(lang.name)
                        : addLang(lang.name)
                    "
                  >
                    <font-awesome-icon
                      v-if="!hasLang(lang.name)"
                      icon="plus"
                    ></font-awesome-icon>
                    <font-awesome-icon
                      v-if="hasLang(lang.name)"
                      icon="minus"
                    ></font-awesome-icon>
                    {{
                      lang.name + " (" + lang.numDocumentsInLanguage + ")"
                    }}</a
                  >
                </div>
              </div>
            </div>
          </div>

          <!-- 文件列表 -->
          <div
            v-if="searchType == 'files' && files != null && files.length > 0"
          >
            <div class="files-list-title">文件列表</div>
            <div class="files-list">
              <div
                class="file-list-item"
                v-for="file in files"
                :key="file.documentID"
              >
                <div v-if="file.filename != null && file.filename != ''">
                  <a
                    href="#"
                    @click="viewFile(file.documentID)"
                    :title="file.filename"
                    >{{ shorten(file.filename) }}</a
                  >
                </div>
              </div>
            </div>
          </div>
        </el-col>

        <!-- 右侧的搜索结果内容展示 -->
        <el-col :span="16">
          <!-- 多个仓库的聚合页面 -->
          <div v-if="!showContent && hasLine" class="results">
            <div class="results-repo" v-for="repo in repos" :key="repo.repoID">
              <!-- 仓库标题 -->
              <div
                class="results-repo-title"
                v-if="
                  repo.repoID != 0 &&
                  repo.documents != null &&
                  repo.documents.length > 0
                "
              >
                仓库：<a
                  href="#"
                  @click="
                    addRepo(
                      repo.remoteURL == '' ? repo.localPath : repo.remoteURL
                    )
                  "
                  >{{
                    simplifyRepo(
                      repo.remoteURL == "" ? repo.localPath : repo.remoteURL
                    )
                  }}</a
                >
                <span
                  class="results-repo-seemore"
                  v-if="repo.documents != null && repos.length > 1"
                  @click="
                    addRepo(
                      repo.remoteURL == '' ? repo.localPath : repo.remoteURL
                    )
                  "
                >
                  该仓库有{{ repo.numDocumentsInRepo }}个匹配的文件
                </span>
              </div>

              <!-- 文件详情 -->
              <div
                :class="doc.class"
                v-for="doc in repo.documents"
                :key="doc.documentID"
              >
                <div>
                  <div class="results-filename">
                    <span
                      >文件：<a href="#" @click="viewFile(doc.documentID)">{{
                        doc.filename
                      }}</a>
                    </span>
                    <span
                      class="results-language"
                      @click="addLang(doc.language)"
                    >
                      {{ doc.language }}
                    </span>
                  </div>
                  <div
                    v-if="doc.lines != null && doc.lines.length > 0"
                    class="results-lines"
                  >
                    <div v-for="line in doc.lines" :key="line.lineNumber">
                      <div :class="line.class">
                        <pre class="inline-pre"><span class="noselect">{{
                    line.lineNumber + 1 + ":\t"
                  }}</span><div v-html="smartEscape(line.content)"></div></pre>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 单文件结果页面 -->
          <div class="results" v-if="showContent && oneRepo != null && hasLine">
            <div class="back" @click="back">
              <font-awesome-icon icon="arrow-left"></font-awesome-icon> 返回
            </div>

            <div class="results-repo">
              <!-- 仓库标题 -->
              <div class="results-repo-title" v-if="oneRepo.repoID != 0">
                仓库：{{
                  simplifyRepo(
                    oneRepo.remoteURL == ""
                      ? oneRepo.localPath
                      : oneRepo.remoteURL
                  )
                }}
              </div>

              <!-- 文件详情 -->
              <div v-for="doc in oneRepo.documents" :key="doc.documentID">
                <div v-if="doc.lines != null && doc.lines.length > 0">
                  <div class="results-filename">
                    <span
                      >文件：<a href="#" @click="viewFile(doc.documentID)">{{
                        doc.filename
                      }}</a>
                    </span>
                    <span
                      class="results-language"
                      @click="addLang(doc.language)"
                    >
                      {{ doc.language }}
                    </span>
                  </div>
                  <div class="results-lines">
                    <div v-for="line in doc.lines" :key="line.lineNumber">
                      <div :class="line.class">
                        <pre class="inline-pre"><span class="noselect">{{
                    line.lineNumber + 1 + ":\t"
                  }}</span><div v-html="smartEscape(line.content)"></div></pre>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </el-main>
  </el-container>
</template>

<script>
import axios from "axios";
import { getFullPathWithQuery } from "@/helpers/query.js";
import { ElMessage } from "element-plus";

export default {
  data() {
    return {
      docID: 0,
      query: "",
      repos: [],
      files: [],
      docs: [],
      oneRepo: null,
      numRepos: 0,
      numDocs: 0,
      numResults: 0,
      duration: 0,
      searched: false,
      hasMore: false,
      message: "",
      showContent: false,
      hasLine: false,
      searchType: "",
      langs: [],
    };
  },

  mounted: function () {
    // 搜索框 focus
    this.$refs.searchInput.$el.children[0].focus();

    // 当 query 为空时，从 URL 参数读入 query
    if (this.query == "") {
      let query = this.$route.query.query;
      if (query == null) {
        query = "";
      }
      this.query = query;
      if (this.query != "") {
        this.clickSearchButton(0, 0, 0, 0);
      }
    }
  },

  methods: {
    syncURL() {
      this.$router.push({
        path: this.$route.path,
        query: {
          query: this.query,
        },
      });
    },

    scrollToTop() {
      window.scrollTo(0, 0);
    },

    shorten(line) {
      let fields = line.split("/");
      if (fields.length > 4) {
        return (
          fields.slice(0, 1).join("/") +
          "/.../" +
          fields.slice(fields.length - 1).join("/")
        );
      }
      return line;
    },

    smartEscape(line) {
      line = line.replaceAll(
        '<b class="keywords">',
        "CONTENTHIGHLIGHTSREPLACEMENTSTART"
      );
      line = line.replaceAll("</b>", "CONTENTHIGHLIGHTSREPLACEMENTEND");
      line = line.replaceAll("<", "&lt;");
      line = line.replaceAll(">", "&gt;");
      line = line.replaceAll(
        "CONTENTHIGHLIGHTSREPLACEMENTSTART",
        '<b class="keywords">'
      );
      line = line.replaceAll("CONTENTHIGHLIGHTSREPLACEMENTEND", "</b>");

      return line;
    },

    simplifyRepo(path) {
      let fields = path.split(":");
      if (fields.length == 2) {
        let repo = fields[1];
        repo = repo.replace(/^\/*/g, "");
        return repo;
      }
      return path;
    },

    addRepo(path) {
      path = this.simplifyRepo(path);
      if (!this.query.includes('repo:"' + path + '"')) {
        this.query = this.query + ' repo:"' + path + '"';
      }
      this.clickSearchButton(0, 0, 0, 0);
    },

    hasRepo(path) {
      path = this.simplifyRepo(path);
      return this.query.includes('repo:"' + path + '"');
    },

    removeRepo(path) {
      path = this.simplifyRepo(path);
      if (this.query.includes('repo:"' + path + '"')) {
        this.query = this.query.replaceAll('repo:"' + path + '"', "").trim();
      }
      this.clickSearchButton(0, 0, 0, 0);
    },

    hasLang(lang) {
      return this.query.includes("lang:^" + lang + "$");
    },

    addLang(lang) {
      if (!this.hasLang(lang)) {
        this.query = this.query + " lang:^" + lang + "$";
      }
      this.clickSearchButton(0, 0, 0, 0);
    },

    removeLang(lang) {
      if (this.hasLang(lang)) {
        this.query = this.query.replaceAll("lang:^" + lang + "$", "").trim();
      }
      this.clickSearchButton(0, 0, 0, 0);
    },

    viewFile(id) {
      this.clickSearchButton(id, 0, 10000, 0);
      this.showContent = true;
    },

    back() {
      this.showContent = false;
    },

    enterSearch(e) {
      if (e.keyCode == 13 && this.query != "") {
        this.clickSearchButton(0, 0, 0, 0);
      }
    },

    addLineSeperator(repos) {
      if (repos == null) {
        return;
      }
      for (let i = 0; i < repos.length; i++) {
        if (repos[i].documents == null) {
          continue;
        }
        for (let j = 0; j < repos[i].documents.length; j++) {
          let doc = repos[i].documents[j];
          if (doc == null || doc.lines == null) {
            continue;
          }
          for (let k = 0; k < doc.lines.length; k++) {
            if (
              k + 1 < doc.lines.length &&
              doc.lines[k].lineNumber + 1 < doc.lines[k + 1].lineNumber
            ) {
              doc.lines[k].class = "line-with-seperator";
            }
          }
        }
      }
    },

    async clickSearchButton(id, r, lc, d) {
      this.docID = id;
      this.searched = true;
      if (this.docID == 0) {
        this.message = "搜索中，请稍等 ...";
        this.showContent = false;
        this.syncURL();
      }

      let searchURL = getFullPathWithQuery("/api/search", {
        q: this.query,
        id: id,
        r: r,
        lc: lc,
        d: d,
      });
      await axios.get(searchURL).then((response) => {
        if (response.data.code == 0) {
          if (this.docID == 0) {
            this.repos = response.data.repos;
            this.searchType = response.data.responseType;

            this.langs = response.data.languages;

            this.addLineSeperator(this.repos);

            // 文件列表
            this.files = [];
            let hasLine = false;
            if (response.data.repos != null) {
              for (let i = 0; i < response.data.repos.length; i++) {
                let repo = response.data.repos[i];
                if (repo.documents != null) {
                  for (let j = 0; j < repo.documents.length; j++) {
                    if (
                      repo.documents[j].lines != null &&
                      repo.documents[j].lines.length > 0
                    ) {
                      hasLine = true;
                    }
                    this.files.push({
                      documentID: repo.documents[j].documentID,
                      filename: repo.documents[j].filename,
                    });
                  }
                }
              }
            }
            this.hasLine = hasLine;

            this.numRepos = response.data.numRepos;
            this.numDocs = response.data.numDocuments;
            this.numResults = response.data.numSections;

            this.duration = response.data.recallDurationInMicroSeconds / 1000;

            if (this.files != null && this.numDocs > this.files.length) {
              this.hasMore = true;
            } else {
              this.hasMore = false;
            }
            this.message =
              "搜索到" +
              this.numRepos +
              "个仓库" +
              this.numDocs +
              "个文件的" +
              this.numResults +
              "条结果，召回耗时" +
              this.duration +
              "毫秒" +
              (this.hasMore
                ? "，只返回前" +
                  this.files.length +
                  "个文件" +
                  (this.repos != null && this.repos.length > 1
                    ? "（请在下方仓库列表中点击缩小范围）"
                    : "（请优化搜索表达式缩小范围）")
                : "");
          } else {
            if (
              response.data.repos != null &&
              response.data.repos.length == 1
            ) {
              this.oneRepo = response.data.repos[0];
            }
          }
        } else {
          ElMessage.error({
            message: response.data.message,
            type: "error",
          });
        }
      });
    },
  },
};
</script>

<style scoped>
dt {
  margin-top: 20px;
  margin-bottom: 5px;
  margin-right: 20px;
  color: #337ab7;
  font-size: 16px;
}
dd {
  margin-right: 20px;
}

.tips {
  margin-top: 20px;
}

.search-results {
  padding-top: 30px;
}

.repos-list {
  overflow-y: scroll;
  height: 400px;
  margin-bottom: 30px;
  background-color: #337ab711;
  padding: 10px;

  border-top: solid;
  border-top-width: 1px;
  padding-left: 10px;
  border-top-color: #bbb;
}
.repos-list-title {
  margin-bottom: 10px;
}
.repo-list-item {
  white-space: pre-wrap;
  word-wrap: break-word;
  padding-top: 2px;
  padding-bottom: 2px;
  margin: 0px;
}
.repo-list-item a {
  color: #337ab7;
}

.files-list {
  overflow-y: scroll;
  max-height: calc(100vh - 700px);
  background-color: #337ab711;
  padding: 10px;

  border-top: solid;
  border-top-width: 1px;
  padding-left: 10px;
  border-top-color: #bbb;
}
.files-list-title {
  margin-bottom: 10px;
}
.file-list-item {
  white-space: pre-wrap;
  word-wrap: break-word;
  padding-top: 2px;
  padding-bottom: 2px;
  margin: 0px;
}
.file-list-item a {
  color: #337ab7;
}

.results {
  overflow-y: scroll;
  max-height: calc(100vh - 170px);
  margin-left: 30px;
}
.results-repo {
  margin-bottom: 50px;
  border-top: solid;
  border-top-width: 1px;
  border-top-color: #bbb;
}
.results-repo-title {
  font-size: 16px;
  font-weight: 700;

  padding-top: 10px;
  font-size: 16px;
  font-weight: 700;
}
.results-repo-title a {
  color: #337ab7;
}
.results-repo-seemore {
  border-radius: 4px;
  background: rgba(51, 122, 183, 0.9);
  color: white;
  padding-top: 2px;
  padding-bottom: 2px;
  padding-left: 6px;
  padding-right: 6px;
  font-size: 12px;
  margin-left: 10px;
  cursor: pointer;
  white-space: nowrap;
}
.results-filename {
  padding-bottom: 10px;
  word-break: break-all;
  word-wrap: break-word;
  white-space: pre-wrap;
  padding-top: 10px;
}
.results-filename a {
  color: #337ab7;
}
.results-language {
  border-radius: 4px;
  background: rgba(51, 122, 183, 0.9);
  color: white;
  padding-top: 2px;
  padding-bottom: 2px;
  padding-left: 6px;
  padding-right: 6px;
  font-size: 12px;
  margin-left: 10px;
  cursor: pointer;
  white-space: nowrap;
}

.back {
  cursor: pointer;
  margin-bottom: 10px;
}

.results-lines {
  background-color: rgba(238, 238, 255, 0.6);
}
.inline-pre {
  white-space: pre-wrap;
  padding-top: 2px;
  padding-bottom: 2px;
  margin: 0px;
  display: flex;
  font-family: Menlo, Monaco, Consolas, "Courier New", monospace;
  word-break: break-all;
  word-wrap: break-word;
}
.noselect {
  user-select: none;
}
.line-with-seperator {
  border-bottom: solid;
  border-bottom-color: lightgray;
  border-bottom-width: 1px;
  padding-bottom: 2px;
  margin-bottom: 2px;
}
</style>
