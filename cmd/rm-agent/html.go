package main

const staticHTML = `<!DOCTYPE html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>

/* simple tree */
.simple-tree {
  user-select: none;
  -moz-user-select: none;
}
.simple-tree>details>summary {
  display: none;
}
.simple-tree a,
.simple-tree summary {
  display: block;
  width: fit-content;
  width: -moz-fit-content;
  border: solid 1px transparent;
  padding: 0 2px;
  outline: none;
  cursor: pointer;
}
.simple-tree a {
  text-decoration: none;
  color: inherit;
}
.simple-tree ::-webkit-details-marker {
  display: none;
}
.simple-tree summary {
  list-style-type: none;
  background-color: #eee;
  outline: none;
}
.simple-tree.dark summary {
  background-color: #444;
}
.simple-tree details>:not(details),
.simple-tree details {
  position: relative;
}
.simple-tree details :not(summary) {
  margin-left: 20px;
}
.simple-tree.nodots details :not(summary) {
  margin-left: 12px;
}
.simple-tree details::before,
.simple-tree details>:not(details)::before {
  content: '';
  width: 10px;
  display: block;
  position: absolute;
}
.simple-tree details::before,
.simple-tree details>:not(details)::before {
  background: url('data:image/svg+xml;utf8,<svg viewBox="0 0 2 2" xmlns="http://www.w3.org/2000/svg"><g><rect x="0" y="0" width="1" height="1"/></g></svg>') left top / 2px 2px;
}
.simple-tree.dark details::before,
.simple-tree.dark details>:not(summary)::before {
  background-image: url('data:image/svg+xml;utf8,<svg viewBox="0 0 2 2" xmlns="http://www.w3.org/2000/svg"><g><rect x="0" y="0" width="1" height="1" fill="white"/></g></svg>');
}
.simple-tree.nodots details::before,
.simple-tree.nodots details>:not(summary)::before {
  background-image: none;
}
.simple-tree details::before {
  top: 0;
  height: 100%;
  background-repeat: repeat-y;
  left: 5px;
  z-index: -1;
}
.simple-tree details>:not(details)::before {
  top: 8px;
  height: calc(100% - 8px);
  background-repeat: repeat-x;
  left: -15px;
}
.simple-tree details>summary::before {
  background: url('data:image/svg+xml;utf8,<svg viewBox="0 0 12 12" xmlns="http://www.w3.org/2000/svg"><g><rect x="0" y="0" width="12" height="12" fill="white" stroke="gray" stroke-width="1"/><line x1="3" y1="6" x2="9" y2="6" stroke="black" stroke-width="2"/><line x1="6" y1="3" x2="6" y2="9" stroke="black" stroke-width="2"/></g></svg>') center center / 12px 12px no-repeat;
  left: -22px;
  top: 2px;
  width: 16px;
  height: 16px;
}
.simple-tree details[open]>summary::before {
  background-image: url('data:image/svg+xml;utf8,<svg viewBox="0 0 12 12" xmlns="http://www.w3.org/2000/svg"><title/><g><rect x="0" y="0" width="12" height="12" fill="white" stroke="gray" stroke-width="1"/><line x1="3" y1="6" x2="9" y2="6" stroke="black" stroke-width="2"/></g></svg>');
}
/* async tree */
.async-tree details[open][data-loaded=false] {
  pointer-events: none;
}
.async-tree details[open][data-loaded=false]>summary::before {
  background-image: url('data:image/svg+xml;utf8,<svg viewBox="0 0 64 64" xmlns="http://www.w3.org/2000/svg"><g><animateTransform attributeName="transform" type="rotate" from="0 32 32" to="360 32 32" dur="1s" repeatCount="indefinite"/><circle cx="32" cy="32" r="32" fill="whitesmoke"/><path d="M 62 32 A 30 30 0 0 0 32 2" style="stroke: black; stroke-width:6; fill:none;"/></g></svg>');
}
.async-tree.black details[open][data-loaded=false]>summary::before {
  background-image: url('data:image/svg+xml;utf8,<svg viewBox="0 0 64 64" xmlns="http://www.w3.org/2000/svg"><g><animateTransform attributeName="transform" type="rotate" from="0 32 32" to="360 32 32" dur="1s" repeatCount="indefinite"/><circle cx="32" cy="32" r="32" fill="whitesmoke"/><path d="M 62 32 A 30 30 0 0 0 32 2" style="stroke: white; stroke-width:6; fill:none;"/></g></svg>');
}
/* select tree */
.select-tree .selected {
  background-color: #beebff;
  border-color: #99defd;
  z-index: 1;
}

.select-tree.dark .selected {
  background-color: #3a484e;
  border-color: #99defd;
}

body {font-family: Arial, Helvetica, sans-serif;}
* {box-sizing: border-box;}

.form-inline {  
  display: flex;
  flex-flow: row wrap;
  align-items: center;
}

.form-inline label {
  margin: 5px 10px 5px 0;
}

.form-inline input {
  vertical-align: middle;
  margin: 5px 10px 5px 0;
  padding: 10px;
  background-color: #fff;
  border: 1px solid #ddd;
}

.form-inline button {
  padding: 10px 20px;
  background-color: dodgerblue;
  border: 1px solid #ddd;
  color: white;
  cursor: pointer;
}

.form-inline button:hover {
  background-color: royalblue;
}

@media (max-width: 800px) {
  .form-inline input {
	margin: 10px 0;
  }
  
  .form-inline {
	flex-direction: column;
	align-items: stretch;
  }
}
</style>
</head>
<body>

<script type="text/javascript" >

'use strict';

{
  const Emitter = typeof window.Emitter === 'undefined' ? class Emitter {
    constructor() {
      this.events = {};
    }
    on(name, callback) {
      this.events[name] = this.events[name] || [];
      this.events[name].push(callback);
    }
    once(name, callback) {
      callback.once = true;
      this.on(name, callback);
    }
    emit(name, ...data) {
      if (this.events[name] === undefined ) {
        return;
      }
      for (const c of [...this.events[name]]) {
        c(...data);
        if (c.once) {
          const index = this.events[name].indexOf(c);
          this.events[name].splice(index, 1);
        }
      }
    }
  } : window.Emitter;

  class SimpleTree extends Emitter {
    constructor(parent, properties = {}) {
      super();
      // do not toggle with click
      parent.addEventListener('click', e => {
        // e.clientX to prevent stopping Enter key
        // e.detail to prevent dbl-click
        // e.offsetX to allow plus and minus clicking
        if (e && e.clientX && e.detail === 1 && e.offsetX >= 0) {
          return e.preventDefault();
        }
        const active = this.active();
        if (active && active.dataset.type === SimpleTree.FILE) {
          e.preventDefault();
          this.emit('action', active);
          if (properties['no-focus-on-action'] === true) {
            window.clearTimeout(this.id);
          }
        }
      });
      parent.classList.add('simple-tree');
      if (properties.dark) {
        parent.classList.add('dark');
      }
      this.parent = parent.appendChild(document.createElement('details'));
      this.parent.appendChild(document.createElement('summary'));
      this.parent.open = true;
      // use this function to alter a node before being passed to this.file or this.folder
      this.interrupt = node => node;
    }
    append(element, parent, before, callback = () => {}) {
      if (before) {
        parent.insertBefore(element, before);
      }
      else {
        parent.appendChild(element);
      }
      callback();
      return element;
    }
    file(node, parent = this.parent, before) {
      parent = parent.closest('details');
      node = this.interrupt(node);
      const a = this.append(Object.assign(document.createElement('a'), {
        textContent: node.name,
        href: '/download/' + node.id
      }), parent, before);
      a.dataset.type = SimpleTree.FILE;
      this.emit('created', a, node);
      return a;
    }
    folder(node, parent = this.parent, before) {
      parent = parent.closest('details');
      node = this.interrupt(node);
      const details = document.createElement('details');
      const summary = Object.assign(document.createElement('summary'), {
        textContent: node.name
      });
      details.appendChild(summary);
      this.append(details, parent, before, () => {
        details.open = node.open;
        details.dataset.type = SimpleTree.FOLDER;
      });
      this.emit('created', summary, node);
      return summary;
    }
    open(details) {
      details.open = true;
    }
    hierarchy(element = this.active()) {
      if (this.parent.contains(element)) {
        const list = [];
        while (element !== this.parent) {
          if (element.dataset.type === SimpleTree.FILE) {
            list.push(element);
          }
          else if (element.dataset.type === SimpleTree.FOLDER) {
            list.push(element.querySelector('summary'));
          }
          element = element.parentElement;
        }
        return list;
      }
      else {
        return [];
      }
    }
    siblings(element = this.parent.querySelector('a, details')) {
      if (this.parent.contains(element)) {
        if (element.dataset.type === undefined) {
          element = element.parentElement;
        }
        return [...element.parentElement.children].filter(e => {
          return e.dataset.type === SimpleTree.FILE || e.dataset.type === SimpleTree.FOLDER;
        }).map(e => {
          if (e.dataset.type === SimpleTree.FILE) {
            return e;
          }
          else {
            return e.querySelector('summary');
          }
        });
      }
      else {
        return [];
      }
    }
    children(details) {
      const e = details.querySelector('a, details');
      if (e) {
        return this.siblings(e);
      }
      else {
        return [];
      }
    }
  }
  SimpleTree.FILE = 'file';
  SimpleTree.FOLDER = 'folder';

  class AsyncTree extends SimpleTree {
    constructor(parent, options) {
      super(parent, options);
      // do not allow toggling when folder is loading
      parent.addEventListener('click', e => {
        const details = e.target.parentElement;
        if (details.open && details.dataset.loaded === 'false') {
          e.preventDefault();
        }
      });
      parent.classList.add('async-tree');
    }
    // add open event for folder creation
    folder(...args) {
      const summary = super.folder(...args);
      const details = summary.closest('details');
      details.addEventListener('toggle', e => {
        this.emit(details.dataset.loaded === 'false' && details.open ? 'fetch' : 'open', summary);
      });
      summary.resolve = () => {
        details.dataset.loaded = true;
        this.emit('open', summary);
      };
      return summary;
    }
    asyncFolder(node, parent, before) {
      const summary = this.folder(node, parent, before);
      const details = summary.closest('details');
      details.dataset.loaded = false;

      if (node.open) {
        this.open(details);
      }

      return summary;
    }
    unloadFolder(summary) {
      const details = summary.closest('details');
      details.open = false;
      const focused = this.active();
      if (focused && this.parent.contains(focused)) {
        this.select(details);
      }
      [...details.children].slice(1).forEach(e => e.remove());
      details.dataset.loaded = false;
    }
    browse(validate, es = this.siblings()) {
      for (const e of es) {
        if (validate(e)) {
          this.select(e);
          if (e.dataset.type === SimpleTree.FILE) {
            return this.emit('browse', e);
          }
          const parent = e.closest('details');
          if (parent.open) {
            return this.browse(validate, this.children(parent));
          }
          else {
            window.setTimeout(() => {
              this.once('open', () => this.browse(validate, this.children(parent)));
              this.open(parent);
            }, 0);
            return;
          }
        }
      }
      this.emit('browse', false);
    }
  }

  class SelectTree extends AsyncTree {
    constructor(parent, options = {}) {
      super(parent, options);
      /* multiple clicks outside of elements */
      parent.addEventListener('click', e => {
        if (e.detail > 1) {
          const active = this.active();
          if (active && active !== e.target) {
            if (e.target.tagName === 'A' || e.target.tagName === 'SUMMARY') {
              return this.select(e.target, 'click');
            }
          }
          if (active) {
            this.focus(active);
          }
        }
      });
      window.addEventListener('focus', () => {
        const active = this.active();
        if (active) {
          this.focus(active);
        }
      });
      parent.addEventListener('focusin', e => {
        const active = this.active();
        if (active !== e.target) {
          this.select(e.target, 'focus');
        }
      });
      this.on('created', (element, node) => {
        if (node.selected) {
          this.select(element);
        }
      });
      parent.classList.add('select-tree');
      // navigate
      if (options.navigate) {
        this.parent.addEventListener('keydown', e => {
          const {code} = e;
          if (code === 'ArrowUp' || code === 'ArrowDown') {
            this.navigate(code === 'ArrowUp' ? 'backward' : 'forward');
            e.preventDefault();
          }
        });
      }
    }
    focus(target) {
      window.clearTimeout(this.id);
      this.id = window.setTimeout(() => document.hasFocus() && target.focus(), 100);
    }
    select(target) {
      const summary = target.querySelector('summary');
      if (summary) {
        target = summary;
      }
      [...this.parent.querySelectorAll('.selected')].forEach(e => e.classList.remove('selected'));
      target.classList.add('selected');
      this.focus(target);
      this.emit('select', target);
    }
    active() {
      return this.parent.querySelector('.selected');
    }
    navigate(direction = 'forward') {
      const e = this.active();
      if (e) {
        const list = [...this.parent.querySelectorAll('a, summary')];
        const index = list.indexOf(e);
        const candidates = direction === 'forward' ? list.slice(index + 1) : list.slice(0, index).reverse();
        for (const m of candidates) {
          if (m.getBoundingClientRect().height) {
            return this.select(m);
          }
        }
      }
    }
  }

  class JSONTree extends SelectTree {
    json(array, parent) {
      array.forEach(item => {
        if (item.type === SimpleTree.FOLDER) {
          const folder = this[item.asynced ? 'asyncFolder' : 'folder'](item, parent);
          if (item.children) {
            this.json(item.children, folder);
          }
        }
        else {
          this.file(item, parent);
        }
      });
    }
  }

  window.Tree = JSONTree;
}

const promiseOfTreeData = fetch("/tree").then(r=>r.json()).then(data => {
  return data;
});
window.onload = async () => {
  let treeData = await promiseOfTreeData;
  var tree = new Tree(document.getElementById('tree'));

  tree.on('action', e => console.log('action', e));
  tree.json(treeData);
  };
</script>   

<h2>Upload File(s)</h2>
<form class="form-inline" action="/upload" method="post" enctype="multipart/form-data" accept=".pdf,.epub">
  <input type="file" id="file" name="file" multiple>
  <button type="submit">Submit</button>
</form>

<div id="tree"></div>

</body>
</html>
`
