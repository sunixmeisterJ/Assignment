'use strict';
let state = {items:[],people:[],rentals:[]};
const $ = id => document.getElementById(id);
const money = n => new Intl.NumberFormat('de-AT',{style:'currency',currency:'EUR'}).format(n/100);
const date = s => new Date(s).toLocaleString('de-AT');
const statusLabel = {reserved:'Reserviert',active:'Ausgeliehen',returned:'Zurückgegeben',cancelled:'Storniert'};
function message(text,error=false){$('message').textContent=text;$('message').className=error?'error':'';}
async function api(path, data){const response=await fetch('/api/'+path,data===undefined?{}:{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(data)});const result=await response.json();if(!response.ok)throw new Error(result.error||'Anfrage fehlgeschlagen');return result;}
function el(tag,text){const node=document.createElement(tag);if(text!==undefined)node.textContent=text;return node;}
function cell(row,text){const c=el('td',text);row.append(c);return c;}
function button(parent,text,fn){const b=el('button',text);b.type='button';b.addEventListener('click',fn);parent.append(b);}
function options(id,rows){const select=$(id),old=select.value;select.replaceChildren();const placeholder=el('option',rows.length?'Bitte wählen':'Zuerst Eintrag anlegen');placeholder.value='';select.append(placeholder);for(const r of rows){const opt=el('option',r.name);opt.value=r.id;select.append(opt);}select.value=old;}
function render(){
const now=Date.now();const busy=new Set(state.rentals.filter(r=>r.status==='active'||(r.status==='reserved'&&Date.parse(r.start)<=now&&Date.parse(r.end)>now)).map(r=>r.itemId));
$('stats').replaceChildren();for(const [label,n] of [['Verfügbar',state.items.length-busy.size],['Ausgeliehen',state.rentals.filter(r=>r.status==='active').length],['Gesamtbestand',state.items.length],['Registrierte Personen',state.people.length]]){const c=el('div',label);c.className='stat';c.append(el('b',n));$('stats').append(c);}
$('items').replaceChildren();for(const i of state.items){const row=el('tr');cell(row,i.name);cell(row,i.details);cell(row,money(i.hourlyCents));button(cell(row),'Bearbeiten',()=>{const f=$('itemForm');f.elements.id.value=i.id;f.elements.name.value=i.name;f.elements.details.value=i.details;f.elements.price.value=(i.hourlyCents/100).toFixed(2);f.elements.name.focus();});$('items').append(row);}
$('persons').replaceChildren();for(const p of state.people){const row=el('tr');cell(row,p.name);cell(row,p.documentType);cell(row,p.documentNumber);$('persons').append(row);}
options('itemSelect',state.items);options('personSelect',state.people);
$('bookings').replaceChildren();for(const r of [...state.rentals].reverse()){const row=el('tr');cell(row,(state.items.find(i=>i.id===r.itemId)?.name||'?')+' / '+(state.people.find(p=>p.id===r.personId)?.name||'?'));cell(row,date(r.start)+' – '+date(r.end)+(r.returnedAt?' · Rückgabe: '+date(r.returnedAt):''));const badge=el('span',statusLabel[r.status]);badge.className='badge';cell(row).append(badge);cell(row,money(r.priceCents)+(r.status==='reserved'||r.status==='active'?' (geplant)':''));const actions=cell(row);if(r.status==='reserved'){if(Date.parse(r.start)<=now&&now<Date.parse(r.end))button(actions,'Ausgeben',()=>transition(r.id,'checkout'));button(actions,'Stornieren',()=>transition(r.id,'cancel'));}if(r.status==='active')button(actions,'Zurückgeben',()=>transition(r.id,'return'));$('bookings').append(row);}
}
async function refresh(){state=await api('state');render();}
async function transition(id,action){try{await api('rentals/'+id+'/'+action,{});await refresh();message('Buchung aktualisiert.');}catch(e){message(e.message,true);}}
function submit(id,path,values){$(id).addEventListener('submit',async e=>{e.preventDefault();const form=e.currentTarget;const b=form.querySelector('button');b.disabled=true;try{await api(path,values(form));form.reset();if(id==='rentalForm')setDates();await refresh();message('Gespeichert.');}catch(err){message(err.message,true);}finally{b.disabled=false;}});}
submit('itemForm','items',f=>({id:Number(f.elements.id.value),name:f.elements.name.value,details:f.elements.details.value,hourlyCents:Math.round(Number(f.elements.price.value)*100)}));
submit('personForm','people',f=>({name:f.elements.name.value,documentType:f.elements.documentType.value,documentNumber:f.elements.documentNumber.value}));
submit('rentalForm','rentals',f=>({itemId:Number(f.elements.itemId.value),personId:Number(f.elements.personId.value),start:new Date(f.elements.start.value).toISOString(),end:new Date(f.elements.end.value).toISOString()}));
$('resetItem').addEventListener('click',()=>$('itemForm').reset());
function localDate(d){return new Date(d.getTime()-d.getTimezoneOffset()*60000).toISOString().slice(0,16);}
function setDates(){const now=new Date();$('rentalForm').elements.start.value=localDate(now);$('rentalForm').elements.end.value=localDate(new Date(now.getTime()+3600000));}
setDates();refresh().catch(e=>message(e.message,true));setInterval(()=>refresh().catch(e=>message(e.message,true)),30000);
